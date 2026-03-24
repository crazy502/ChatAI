package chat

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"server/infra/mq"
	"server/internal/ai"

	"github.com/streadway/amqp"
)

type MessageQueuePayload struct {
	IdempotencyKey string `json:"idempotency_key"`
	SessionID      string `json:"session_id"`
	Content        string `json:"content"`
	UserName       string `json:"user_name"`
	IsUser         bool   `json:"is_user"`
}

var consumerOnce sync.Once

func marshalMessagePayload(message *ai.StoredMessage) ([]byte, error) {
	payload := MessageQueuePayload{
		IdempotencyKey: message.IdempotencyKey,
		SessionID:      message.SessionID,
		Content:        message.Content,
		UserName:       message.UserName,
		IsUser:         message.IsUser,
	}
	return json.Marshal(payload)
}

func StartMessageConsumer(repo *Repository) {
	if mq.RMQConsumer == nil || repo == nil {
		return
	}

	consumerOnce.Do(func() {
		go mq.RMQConsumer.Consume(func(msg *amqp.Delivery) error {
			return consumeMessage(msg, repo)
		})
	})
}

func consumeMessage(msg *amqp.Delivery, repo *Repository) error {
	var payload MessageQueuePayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		return fmt.Errorf("%w: invalid message body: %v", mq.ErrDropMessage, err)
	}

	if payload.IdempotencyKey == "" {
		payload.IdempotencyKey = legacyMessageID(msg.Body)
	}

	_, err := repo.Create(&Message{
		IdempotencyKey: payload.IdempotencyKey,
		SessionID:      payload.SessionID,
		Content:        payload.Content,
		UserName:       payload.UserName,
		IsUser:         payload.IsUser,
	})
	if err != nil {
		return fmt.Errorf("persist message failed: %w", err)
	}

	return nil
}

func saveWithQueue(repo *Repository, message *ai.StoredMessage) error {
	data, err := marshalMessagePayload(message)
	if err == nil && mq.RMQMessage != nil {
		if publishErr := mq.RMQMessage.Publish(data); publishErr == nil {
			return nil
		} else {
			log.Printf("rabbitmq publish failed, fallback to direct db insert: %v", publishErr)
		}
	}

	_, err = repo.Create(&Message{
		IdempotencyKey: message.IdempotencyKey,
		SessionID:      message.SessionID,
		Content:        message.Content,
		UserName:       message.UserName,
		IsUser:         message.IsUser,
	})
	return err
}

func toAIStoredMessages(messages []Message) []ai.StoredMessage {
	result := make([]ai.StoredMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, ai.StoredMessage{
			IdempotencyKey: message.IdempotencyKey,
			SessionID:      message.SessionID,
			UserName:       message.UserName,
			Content:        message.Content,
			IsUser:         message.IsUser,
		})
	}
	return result
}

func legacyMessageID(body []byte) string {
	sum := sha256.Sum256(body)
	return "legacy-mq-" + hex.EncodeToString(sum[:])
}
