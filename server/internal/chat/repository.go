package chat

import (
	"strings"

	"server/infra/db"
	"server/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetMessagesBySessionID(sessionID string) ([]Message, error) {
	var messages []Message
	err := db.DB.Where("session_id = ?", sessionID).
		Order("created_at asc, id asc").
		Find(&messages).
		Error
	return messages, err
}

func (r *Repository) Create(message *Message) (*Message, error) {
	if strings.TrimSpace(message.IdempotencyKey) == "" {
		message.IdempotencyKey = utils.GenerateUUID()
	}

	result := db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idempotency_key"}},
		DoNothing: true,
	}).Create(message)
	if result.Error != nil {
		return message, result.Error
	}

	if result.RowsAffected > 0 {
		return message, nil
	}

	existing := new(Message)
	err := db.DB.Where("idempotency_key = ?", message.IdempotencyKey).First(existing).Error
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *Repository) GetAll() ([]Message, error) {
	var messages []Message
	err := db.DB.Order("created_at asc, id asc").Find(&messages).Error
	return messages, err
}

func (r *Repository) EnsureMessageIdempotency() error {
	if err := db.DB.Model(&Message{}).
		Where("idempotency_key IS NULL OR idempotency_key = ''").
		Update("idempotency_key", gorm.Expr("CONCAT('legacy-', id)")).Error; err != nil {
		return err
	}

	if db.DB.Migrator().HasIndex(&Message{}, "idx_messages_idempotency_key") {
		return nil
	}

	return db.DB.Exec("CREATE UNIQUE INDEX idx_messages_idempotency_key ON messages (idempotency_key)").Error
}
