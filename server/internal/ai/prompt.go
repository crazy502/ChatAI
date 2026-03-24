package ai

import (
	"server/pkg/utils"

	"github.com/cloudwego/eino/schema"
)

type PromptMessage struct {
	Content string
	IsUser  bool
}

type StoredMessage struct {
	IdempotencyKey string
	SessionID      string
	UserName       string
	Content        string
	IsUser         bool
}

func NewStoredMessage(sessionID, userName, content string, isUser bool) *StoredMessage {
	return &StoredMessage{
		IdempotencyKey: utils.GenerateUUID(),
		SessionID:      sessionID,
		UserName:       userName,
		Content:        content,
		IsUser:         isUser,
	}
}

func ToPromptMessages(messages []StoredMessage) []PromptMessage {
	result := make([]PromptMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, PromptMessage{
			Content: message.Content,
			IsUser:  message.IsUser,
		})
	}
	return result
}

func ToSchemaMessages(messages []PromptMessage) []*schema.Message {
	result := make([]*schema.Message, 0, len(messages))
	for _, message := range messages {
		role := schema.Assistant
		if message.IsUser {
			role = schema.User
		}

		result = append(result, &schema.Message{
			Role:    role,
			Content: message.Content,
		})
	}
	return result
}
