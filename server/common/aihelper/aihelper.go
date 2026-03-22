package aihelper

import (
	"context"
	"log"
	"sync"
	"time"

	"server/common/metrics"
	"server/common/rabbitmq"
	messageDAO "server/dao/message"
	"server/model"
	"server/utils"
)

// AIHelper stores conversation history and the bound AI model.
type AIHelper struct {
	model     AIModel
	messages  []*model.Message
	mu        sync.RWMutex
	SessionID string
	saveFunc  func(*model.Message) (*model.Message, error)
}

func NewAIHelper(model_ AIModel, sessionID string) *AIHelper {
	return &AIHelper{
		model:    model_,
		messages: make([]*model.Message, 0),
		saveFunc: func(msg *model.Message) (*model.Message, error) {
			data := rabbitmq.GenerateMessageMQParam(msg)
			if rabbitmq.RMQMessage != nil {
				if err := rabbitmq.RMQMessage.Publish(data); err == nil {
					return msg, nil
				} else {
					log.Printf("rabbitmq publish failed, fallback to direct db insert: %v", err)
				}
			}

			_, err := messageDAO.CreateMessage(msg)
			return msg, err
		},
		SessionID: sessionID,
	}
}

func (a *AIHelper) AddMessage(content string, userName string, isUser bool, save bool) {
	userMsg := model.Message{
		IdempotencyKey: utils.GenerateUUID(),
		SessionID:      a.SessionID,
		Content:        content,
		UserName:       userName,
		IsUser:         isUser,
	}

	a.mu.Lock()
	a.messages = append(a.messages, &userMsg)
	a.mu.Unlock()

	if save {
		if _, err := a.saveFunc(&userMsg); err != nil {
			log.Printf("persist message failed for session=%s user=%s: %v", a.SessionID, userName, err)
		}
	}
}

func (a *AIHelper) SetSaveFunc(saveFunc func(*model.Message) (*model.Message, error)) {
	a.saveFunc = saveFunc
}

func (a *AIHelper) HasMessages() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.messages) > 0
}

func (a *AIHelper) ReplaceMessages(history []model.Message) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.messages = make([]*model.Message, 0, len(history))
	for i := range history {
		msg := history[i]
		a.messages = append(a.messages, &msg)
	}
}

func (a *AIHelper) GetMessages() []*model.Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]*model.Message, len(a.messages))
	copy(out, a.messages)
	return out
}

func (a *AIHelper) GenerateResponse(userName string, ctx context.Context, userQuestion string) (*model.Message, error) {
	a.AddMessage(userQuestion, userName, true, true)

	a.mu.RLock()
	messages := utils.ConvertToSchemaMessages(a.messages)
	a.mu.RUnlock()

	start := time.Now()
	schemaMsg, err := a.model.GenerateResponse(ctx, messages)
	metrics.GetCollector().RecordModel(a.GetModelType(), "generate", userName, time.Since(start), err)
	if err != nil {
		return nil, err
	}

	modelMsg := utils.ConvertToModelMessage(a.SessionID, userName, schemaMsg)
	a.AddMessage(modelMsg.Content, userName, false, true)

	return modelMsg, nil
}

func (a *AIHelper) StreamResponse(userName string, ctx context.Context, cb StreamCallback, userQuestion string) (*model.Message, error) {
	a.AddMessage(userQuestion, userName, true, true)

	a.mu.RLock()
	messages := utils.ConvertToSchemaMessages(a.messages)
	a.mu.RUnlock()

	start := time.Now()
	content, err := a.model.StreamResponse(ctx, messages, cb)
	metrics.GetCollector().RecordModel(a.GetModelType(), "stream", userName, time.Since(start), err)
	if err != nil {
		return nil, err
	}

	modelMsg := &model.Message{
		SessionID: a.SessionID,
		UserName:  userName,
		Content:   content,
		IsUser:    false,
	}

	a.AddMessage(modelMsg.Content, userName, false, true)
	return modelMsg, nil
}

func (a *AIHelper) GetModelType() string {
	return a.model.GetModelType()
}
