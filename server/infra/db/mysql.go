package db

import (
	"fmt"
	"time"

	"server/infra/config"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMysql() error {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		cfg.MysqlUser,
		cfg.MysqlPassword,
		cfg.MysqlHost,
		cfg.MysqlPort,
		cfg.MysqlDatabaseName,
		cfg.MysqlCharset,
	)

	var gormLogger logger.Interface
	if gin.Mode() == gin.DebugMode {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default
	}

	instance, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	sqlDB, err := instance.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = instance
	return nil
}

func Migrate(models ...any) error {
	return DB.AutoMigrate(models...)
}
