package message

import (
	"strings"

	"server/common/mysql"
	"server/model"
	"server/utils"

	"gorm.io/gorm/clause"
)

func GetMessagesBySessionID(sessionID string) ([]model.Message, error) {
	var msgs []model.Message
	err := mysql.DB.Where("session_id = ?", sessionID).Order("created_at asc, id asc").Find(&msgs).Error
	return msgs, err
}

func GetMessagesBySessionIDs(sessionIDs []string) ([]model.Message, error) {
	var msgs []model.Message
	if len(sessionIDs) == 0 {
		return msgs, nil
	}
	err := mysql.DB.Where("session_id IN ?", sessionIDs).Order("created_at asc, id asc").Find(&msgs).Error
	return msgs, err
}

func CreateMessage(message *model.Message) (*model.Message, error) {
	if strings.TrimSpace(message.IdempotencyKey) == "" {
		message.IdempotencyKey = utils.GenerateUUID()
	}

	result := mysql.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idempotency_key"}},
		DoNothing: true,
	}).Create(message)
	if result.Error != nil {
		return message, result.Error
	}

	if result.RowsAffected > 0 {
		return message, nil
	}

	existing := new(model.Message)
	err := mysql.DB.Where("idempotency_key = ?", message.IdempotencyKey).First(existing).Error
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func GetAllMessages() ([]model.Message, error) {
	var msgs []model.Message
	err := mysql.DB.Order("created_at asc, id asc").Find(&msgs).Error
	return msgs, err
}
