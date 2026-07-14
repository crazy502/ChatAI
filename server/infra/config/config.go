package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"server/pkg/observe"

	"github.com/BurntSushi/toml"
)

type MainConfig struct {
	Port        int    `toml:"port"`
	AppName     string `toml:"appName"`
	Host        string `toml:"host"`
	Environment string `toml:"environment"`
}

type EmailConfig struct {
	Authcode string `toml:"authcode"`
	Email    string `toml:"email"`
}

type RedisConfig struct {
	RedisPort     int    `toml:"port"`
	RedisDB       int    `toml:"db"`
	RedisHost     string `toml:"host"`
	RedisPassword string `toml:"password"`
}

type MysqlConfig struct {
	MysqlPort         int                `toml:"port"`
	MysqlHost         string             `toml:"host"`
	MysqlUser         string             `toml:"user"`
	MysqlPassword     string             `toml:"password"`
	MysqlDatabaseName string             `toml:"databaseName"`
	MysqlCharset      string             `toml:"charset"`
	Pool              MysqlPoolConfig    `toml:"pool"`
	Replica           MysqlReplicaConfig `toml:"replica"`
}

type MysqlPoolConfig struct {
	MaxIdleConns           int `toml:"maxIdleConns"`
	MaxOpenConns           int `toml:"maxOpenConns"`
	ConnMaxLifetimeMinutes int `toml:"connMaxLifetimeMinutes"`
	ConnMaxIdleTimeMinutes int `toml:"connMaxIdleTimeMinutes"`
}

type MysqlReplicaConfig struct {
	Enabled      bool            `toml:"enabled"`
	Host         string          `toml:"host"`
	Port         int             `toml:"port"`
	User         string          `toml:"user"`
	Password     string          `toml:"password"`
	DatabaseName string          `toml:"databaseName"`
	Charset      string          `toml:"charset"`
	Pool         MysqlPoolConfig `toml:"pool"`
}

type JWTConfig struct {
	ExpireDuration int    `toml:"expire_duration"`
	Issuer         string `toml:"issuer"`
	Subject        string `toml:"subject"`
	Key            string `toml:"key"`
}

type AdminConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Email    string `toml:"email"`
}

type RabbitMQConfig struct {
	RabbitMQPort     int    `toml:"port"`
	RabbitMQHost     string `toml:"host"`
	RabbitMQUsername string `toml:"username"`
	RabbitMQPassword string `toml:"password"`
	RabbitMQVhost    string `toml:"vhost"`
}

type QwenConfig struct {
	APIKey    string `toml:"apiKey"`
	BaseURL   string `toml:"baseURL"`
	ModelName string `toml:"modelName"`
}

type DeepSeekConfig struct {
	APIKey    string `toml:"apiKey"`
	BaseURL   string `toml:"baseURL"`
	ModelName string `toml:"modelName"`
}

type Config struct {
	EmailConfig    `toml:"emailConfig"`
	RedisConfig    `toml:"redisConfig"`
	MysqlConfig    `toml:"mysqlConfig"`
	JWTConfig      `toml:"jwtConfig"`
	AdminConfig    AdminConfig `toml:"adminConfig"`
	MainConfig     `toml:"mainConfig"`
	RabbitMQConfig `toml:"rabbitmqConfig"`
	QwenConfig     `toml:"qwenConfig"`
	DeepSeekConfig `toml:"deepseekConfig"`
}

type RedisKeyConfig struct {
	CaptchaPrefix string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	CaptchaPrefix: "captcha:%s",
}

var (
	cfg   *Config
	cfgMu sync.RWMutex
)

func InitConfig() error {
	loaded := new(Config)
	configPath := strings.TrimSpace(os.Getenv("APP_CONFIG_PATH"))
	if configPath == "" {
		configPath = "config/config.toml"
	}

	if _, err := toml.DecodeFile(configPath, loaded); err != nil {
		observe.Error(context.Background(), "加载配置文件失败", err)
		return fmt.Errorf("load config file %q: %w", configPath, err)
	}
	if err := applyEnvironment(loaded); err != nil {
		return err
	}
	applyDefaults(loaded)
	if err := loaded.Validate(); err != nil {
		return err
	}

	cfgMu.Lock()
	cfg = loaded
	cfgMu.Unlock()
	return nil
}

func GetConfig() *Config {
	cfgMu.RLock()
	current := cfg
	cfgMu.RUnlock()
	if current != nil {
		return current
	}

	if err := InitConfig(); err != nil {
		observe.Error(context.Background(), "load application config failed", err)
		return new(Config)
	}

	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg
}

func (c *Config) IsProduction() bool {
	return c != nil && strings.EqualFold(strings.TrimSpace(c.Environment), "production")
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}
	if strings.TrimSpace(c.Key) == "" {
		return errors.New("JWT secret is required; set JWT_SECRET")
	}
	if strings.TrimSpace(c.AdminConfig.Username) == "" || strings.TrimSpace(c.AdminConfig.Email) == "" {
		return errors.New("admin username and email are required")
	}
	if !c.IsProduction() {
		return nil
	}

	secret := strings.TrimSpace(c.Key)
	if len(secret) < 32 || secret == "replace-with-a-random-jwt-secret" {
		return errors.New("production JWT secret must contain at least 32 random characters")
	}
	adminPassword := strings.TrimSpace(c.AdminConfig.Password)
	if len(adminPassword) < 12 || strings.EqualFold(adminPassword, "admin") {
		return errors.New("production admin password must contain at least 12 characters")
	}
	if strings.TrimSpace(c.Email) == "" || strings.TrimSpace(c.Authcode) == "" {
		return errors.New("production mail credentials are required")
	}
	return nil
}

func applyDefaults(c *Config) {
	if strings.TrimSpace(c.AppName) == "" {
		c.AppName = "AgentGo"
	}
	if strings.TrimSpace(c.Host) == "" {
		c.Host = "0.0.0.0"
	}
	if c.Port <= 0 {
		c.Port = 9090
	}
	if c.ExpireDuration <= 0 {
		c.ExpireDuration = 2
	}
	if strings.TrimSpace(c.Issuer) == "" {
		c.Issuer = c.AppName
	}
	if strings.TrimSpace(c.Subject) == "" {
		c.Subject = c.AppName
	}
	if strings.TrimSpace(c.AdminConfig.Username) == "" {
		c.AdminConfig.Username = "admin@qq.com"
	}
	if strings.TrimSpace(c.AdminConfig.Email) == "" {
		c.AdminConfig.Email = "admin@qq.com"
	}
}

func applyEnvironment(c *Config) error {
	stringOverrides := map[string]*string{
		"APP_ENV":                &c.Environment,
		"APP_HOST":               &c.Host,
		"MAIL_ADDRESS":           &c.Email,
		"MAIL_AUTH_CODE":         &c.Authcode,
		"REDIS_HOST":             &c.RedisHost,
		"REDIS_PASSWORD":         &c.RedisPassword,
		"MYSQL_HOST":             &c.MysqlHost,
		"MYSQL_USER":             &c.MysqlUser,
		"MYSQL_PASSWORD":         &c.MysqlPassword,
		"MYSQL_DATABASE":         &c.MysqlDatabaseName,
		"MYSQL_REPLICA_HOST":     &c.MysqlConfig.Replica.Host,
		"MYSQL_REPLICA_USER":     &c.MysqlConfig.Replica.User,
		"MYSQL_REPLICA_PASSWORD": &c.MysqlConfig.Replica.Password,
		"JWT_SECRET":             &c.Key,
		"JWT_ISSUER":             &c.Issuer,
		"JWT_SUBJECT":            &c.Subject,
		"ADMIN_USERNAME":         &c.AdminConfig.Username,
		"ADMIN_PASSWORD":         &c.AdminConfig.Password,
		"ADMIN_EMAIL":            &c.AdminConfig.Email,
		"RABBITMQ_HOST":          &c.RabbitMQHost,
		"RABBITMQ_USERNAME":      &c.RabbitMQUsername,
		"RABBITMQ_PASSWORD":      &c.RabbitMQPassword,
		"QWEN_API_KEY":           &c.QwenConfig.APIKey,
		"DEEPSEEK_API_KEY":       &c.DeepSeekConfig.APIKey,
	}
	for name, target := range stringOverrides {
		if value, ok := os.LookupEnv(name); ok {
			*target = strings.TrimSpace(value)
		}
	}

	intOverrides := map[string]*int{
		"APP_PORT":         &c.Port,
		"REDIS_PORT":       &c.RedisPort,
		"REDIS_DB":         &c.RedisDB,
		"MYSQL_PORT":       &c.MysqlPort,
		"JWT_EXPIRE_HOURS": &c.ExpireDuration,
		"RABBITMQ_PORT":    &c.RabbitMQPort,
	}
	for name, target := range intOverrides {
		value, ok := os.LookupEnv(name)
		if !ok {
			continue
		}
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse %s: %w", name, err)
		}
		*target = parsed
	}

	if value, ok := os.LookupEnv("MYSQL_REPLICA_ENABLED"); ok {
		parsed, err := strconv.ParseBool(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse MYSQL_REPLICA_ENABLED: %w", err)
		}
		c.MysqlConfig.Replica.Enabled = parsed
	}
	return nil
}
