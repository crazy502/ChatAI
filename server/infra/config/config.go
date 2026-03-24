package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type MainConfig struct {
	Port    int    `toml:"port"`
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
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
	MysqlPort         int    `toml:"port"`
	MysqlHost         string `toml:"host"`
	MysqlUser         string `toml:"user"`
	MysqlPassword     string `toml:"password"`
	MysqlDatabaseName string `toml:"databaseName"`
	MysqlCharset      string `toml:"charset"`
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

var cfg *Config

func InitConfig() error {
	if cfg == nil {
		cfg = new(Config)
	}

	if _, err := toml.DecodeFile("config/config.toml", cfg); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetConfig() *Config {
	if cfg == nil {
		cfg = new(Config)
		_ = InitConfig()
	}
	return cfg
}
