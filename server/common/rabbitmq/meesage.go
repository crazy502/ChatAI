package rabbitmq

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	messageDAO "server/dao/message"
	"server/model"

	"github.com/streadway/amqp"
)

var ErrDropMessage = errors.New("drop message")

func legacyMessageIdempotencyKey(body []byte) string {
	sum := sha256.Sum256(body)
	return "legacy-mq-" + hex.EncodeToString(sum[:])
}

type MessageMQParam struct {
	IdempotencyKey string `json:"idempotency_key"`
	SessionID      string `json:"session_id"`
	Content        string `json:"content"`
	UserName       string `json:"user_name"`
	IsUser         bool   `json:"is_user"`
}

func GenerateMessageMQParam(message *model.Message) []byte {
	param := MessageMQParam{
		IdempotencyKey: message.IdempotencyKey,
		SessionID:      message.SessionID,
		Content:        message.Content,
		UserName:       message.UserName,
		IsUser:         message.IsUser,
	}
	data, _ := json.Marshal(param)
	return data
}

func MQMessage(msg *amqp.Delivery) error {
	var param MessageMQParam
	if err := json.Unmarshal(msg.Body, &param); err != nil {
		return fmt.Errorf("%w: invalid message body: %v", ErrDropMessage, err)
	}

	if param.IdempotencyKey == "" {
		param.IdempotencyKey = legacyMessageIdempotencyKey(msg.Body)
	}

	newMsg := &model.Message{
		IdempotencyKey: param.IdempotencyKey,
		SessionID:      param.SessionID,
		Content:        param.Content,
		UserName:       param.UserName,
		IsUser:         param.IsUser,
	}

	if _, err := messageDAO.CreateMessage(newMsg); err != nil {
		return fmt.Errorf("persist message failed: %w", err)
	}

	return nil
}
