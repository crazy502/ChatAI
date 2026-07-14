package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"server/infra/config"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultMysqlCharset    = "utf8mb4"
	defaultMaxIdleConns    = 10
	defaultMaxOpenConns    = 100
	defaultConnMaxLifetime = time.Hour
	defaultConnMaxIdleTime = 10 * time.Minute
)

var (
	DB     *gorm.DB
	readDB *gorm.DB
)

type mysqlConnectionConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	Charset      string
	Pool         config.MysqlPoolConfig
}

type normalizedPoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func InitMysql() error {
	cfg := config.GetConfig()
	gormLogger := buildGormLogger()

	writer, err := openConnection(buildPrimaryConnectionConfig(cfg), gormLogger)
	if err != nil {
		return err
	}

	DB = writer
	readDB = writer

	if readerCfg, ok := buildReplicaConnectionConfig(cfg); ok {
		reader, err := openConnection(readerCfg, gormLogger)
		if err != nil {
			return err
		}
		readDB = reader
	}

	return nil
}

func Writer() *gorm.DB {
	return DB
}

func Reader() *gorm.DB {
	if readDB != nil {
		return readDB
	}
	return DB
}

func HasReplica() bool {
	return DB != nil && readDB != nil && DB != readDB
}

func Close() error {
	var firstErr error
	closed := make(map[*sql.DB]struct{}, 2)
	for _, instance := range []*gorm.DB{readDB, DB} {
		if instance == nil {
			continue
		}
		sqlDB, err := instance.DB()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if _, exists := closed[sqlDB]; exists {
			continue
		}
		closed[sqlDB] = struct{}{}
		if err := sqlDB.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	DB = nil
	readDB = nil
	return firstErr
}

func Migrate(models ...any) error {
	return Writer().AutoMigrate(models...)
}

func buildGormLogger() logger.Interface {
	if gin.Mode() == gin.DebugMode {
		return logger.Default.LogMode(logger.Info)
	}
	return logger.Default
}

func buildPrimaryConnectionConfig(cfg *config.Config) mysqlConnectionConfig {
	return mysqlConnectionConfig{
		Host:         cfg.MysqlHost,
		Port:         cfg.MysqlPort,
		User:         cfg.MysqlUser,
		Password:     cfg.MysqlPassword,
		DatabaseName: cfg.MysqlDatabaseName,
		Charset:      firstNonEmpty(cfg.MysqlCharset, defaultMysqlCharset),
		Pool:         cfg.MysqlConfig.Pool,
	}
}

func buildReplicaConnectionConfig(cfg *config.Config) (mysqlConnectionConfig, bool) {
	replica := cfg.MysqlConfig.Replica
	if !replica.Enabled {
		return mysqlConnectionConfig{}, false
	}

	return mysqlConnectionConfig{
		Host:         firstNonEmpty(replica.Host, cfg.MysqlHost),
		Port:         firstNonZero(replica.Port, cfg.MysqlPort),
		User:         firstNonEmpty(replica.User, cfg.MysqlUser),
		Password:     firstNonEmpty(replica.Password, cfg.MysqlPassword),
		DatabaseName: firstNonEmpty(replica.DatabaseName, cfg.MysqlDatabaseName),
		Charset:      firstNonEmpty(replica.Charset, firstNonEmpty(cfg.MysqlCharset, defaultMysqlCharset)),
		Pool:         mergePoolConfig(cfg.MysqlConfig.Pool, replica.Pool),
	}, true
}

func mergePoolConfig(base, override config.MysqlPoolConfig) config.MysqlPoolConfig {
	merged := base
	if override.MaxIdleConns > 0 {
		merged.MaxIdleConns = override.MaxIdleConns
	}
	if override.MaxOpenConns > 0 {
		merged.MaxOpenConns = override.MaxOpenConns
	}
	if override.ConnMaxLifetimeMinutes > 0 {
		merged.ConnMaxLifetimeMinutes = override.ConnMaxLifetimeMinutes
	}
	if override.ConnMaxIdleTimeMinutes > 0 {
		merged.ConnMaxIdleTimeMinutes = override.ConnMaxIdleTimeMinutes
	}
	return merged
}

func openConnection(cfg mysqlConnectionConfig, gormLogger logger.Interface) (*gorm.DB, error) {
	instance, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       cfg.dsn(),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := instance.DB()
	if err != nil {
		return nil, err
	}
	applyPoolSettings(sqlDB, cfg.Pool)
	return instance, nil
}

func applyPoolSettings(sqlDB *sql.DB, poolCfg config.MysqlPoolConfig) {
	pool := normalizePoolConfig(poolCfg)
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(pool.ConnMaxIdleTime)
}

func normalizePoolConfig(poolCfg config.MysqlPoolConfig) normalizedPoolConfig {
	pool := normalizedPoolConfig{
		MaxIdleConns:    poolCfg.MaxIdleConns,
		MaxOpenConns:    poolCfg.MaxOpenConns,
		ConnMaxLifetime: time.Duration(poolCfg.ConnMaxLifetimeMinutes) * time.Minute,
		ConnMaxIdleTime: time.Duration(poolCfg.ConnMaxIdleTimeMinutes) * time.Minute,
	}

	if pool.MaxIdleConns <= 0 {
		pool.MaxIdleConns = defaultMaxIdleConns
	}
	if pool.MaxOpenConns <= 0 {
		pool.MaxOpenConns = defaultMaxOpenConns
	}
	if pool.ConnMaxLifetime <= 0 {
		pool.ConnMaxLifetime = defaultConnMaxLifetime
	}
	if pool.ConnMaxIdleTime <= 0 {
		pool.ConnMaxIdleTime = defaultConnMaxIdleTime
	}

	return pool
}

func (cfg mysqlConnectionConfig) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
		cfg.Charset,
	)
}

func firstNonEmpty(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value != "" {
		return value
	}
	return strings.TrimSpace(fallback)
}

func firstNonZero(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
