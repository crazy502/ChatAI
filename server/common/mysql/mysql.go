package mysql

import (
	"fmt"
	"log"
	"time"

	"server/config"
	"server/model"
	"server/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "admin"
	defaultAdminEmail    = "admin@gopherai.local"
)

func InitMysql() error {
	conf := config.GetConfig()
	host := conf.MysqlHost
	port := conf.MysqlPort
	dbname := conf.MysqlDatabaseName
	username := conf.MysqlUser
	password := conf.MysqlPassword
	charset := conf.MysqlCharset

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", username, password, host, port, dbname, charset)

	var log logger.Interface
	if gin.Mode() == "debug" {
		log = logger.Default.LogMode(logger.Info)
	} else {
		log = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	return migration()
}

func migration() error {
	if err := DB.AutoMigrate(
		new(model.User),
		new(model.Session),
		new(model.Message),
	); err != nil {
		return err
	}

	if err := migrateMessageIdempotency(); err != nil {
		return err
	}

	return ensureDefaultAdminUser()
}

func migrateMessageIdempotency() error {
	if err := DB.Model(&model.Message{}).
		Where("idempotency_key IS NULL OR idempotency_key = ''").
		Update("idempotency_key", gorm.Expr("CONCAT('legacy-', id)")).Error; err != nil {
		return err
	}

	if DB.Migrator().HasIndex(&model.Message{}, "idx_messages_idempotency_key") {
		return nil
	}

	return DB.Exec("CREATE UNIQUE INDEX idx_messages_idempotency_key ON messages (idempotency_key)").Error
}

func ensureDefaultAdminUser() error {
	passwordHash, err := utils.HashPassword(defaultAdminPassword)
	if err != nil {
		return err
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).
			Where("username <> ?", defaultAdminUsername).
			Update("is_admin", false).Error; err != nil {
			return err
		}

		var adminUser model.User
		err := tx.Where("username = ?", defaultAdminUsername).First(&adminUser).Error
		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(&model.User{
				Email:    defaultAdminEmail,
				Name:     defaultAdminUsername,
				Username: defaultAdminUsername,
				Password: passwordHash,
				IsAdmin:  true,
			}).Error; err != nil {
				return err
			}

			log.Printf("default admin created: username=%s", defaultAdminUsername)
			return nil
		}
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"is_admin": true,
			"name":     defaultAdminUsername,
			"password": passwordHash,
		}
		if adminUser.Email == "" {
			updates["email"] = defaultAdminEmail
		}

		if err := tx.Model(&model.User{}).
			Where("id = ?", adminUser.ID).
			Updates(updates).Error; err != nil {
			return err
		}

		log.Printf("default admin ensured: username=%s", defaultAdminUsername)
		return nil
	})
}

func InsertUser(user *model.User) (*model.User, error) {
	err := DB.Create(&user).Error
	return user, err
}

func GetUserByUsername(username string) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("username = ?", username).First(user).Error
	return user, err
}

func GetUserByID(id int64) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("id = ?", id).First(user).Error
	return user, err
}

func GetUserByEmail(email string) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("email = ?", email).First(user).Error
	return user, err
}
