package chat

import (
	"strings"

	"server/infra/db"
	"server/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct{}

const historyQueryIndexName = "idx_messages_session_created_id"

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetMessagesBySessionID(sessionID string) ([]Message, error) {
	var messages []Message
	err := db.Reader().Where("session_id = ?", sessionID).
		Order("created_at asc, id asc").
		Find(&messages).
		Error
	return messages, err
}

func (r *Repository) Create(message *Message) (*Message, error) {
	if strings.TrimSpace(message.IdempotencyKey) == "" {
		message.IdempotencyKey = utils.GenerateUUID()
	}

	result := db.Writer().Clauses(clause.OnConflict{
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
	err := db.Writer().Where("idempotency_key = ?", message.IdempotencyKey).First(existing).Error
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *Repository) GetAll() ([]Message, error) {
	var messages []Message
	err := db.Reader().Order("created_at asc, id asc").Find(&messages).Error
	return messages, err
}

func (r *Repository) EnsureMessageIdempotency() error {
	if err := db.Writer().Model(&Message{}).
		Where("idempotency_key IS NULL OR idempotency_key = ''").
		Update("idempotency_key", gorm.Expr("CONCAT('legacy-', id)")).Error; err != nil {
		return err
	}

	if db.Writer().Migrator().HasIndex(&Message{}, "idx_messages_idempotency_key") {
		return nil
	}

	return db.Writer().Exec("CREATE UNIQUE INDEX idx_messages_idempotency_key ON messages (idempotency_key)").Error
}

func (r *Repository) EnsureHistoryIndexes() error {
	writer := db.Writer()
	if writer.Migrator().HasIndex(&Message{}, historyQueryIndexName) {
		return nil
	}

	return writer.Exec(
		"CREATE INDEX " + historyQueryIndexName + " ON messages (session_id, created_at, id)",
	).Error
}
