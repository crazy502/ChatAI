package ai

import (
	"context"
	"sync"
	"time"

	"server/infra/metrics"
)

type SaveFunc func(*StoredMessage) error

type Helper struct {
	provider  Provider
	messages  []PromptMessage
	mu        sync.RWMutex
	SessionID string
	saveFunc  SaveFunc
}

func NewHelper(provider Provider, sessionID string) *Helper {
	return &Helper{
		provider:  provider,
		messages:  make([]PromptMessage, 0),
		SessionID: sessionID,
	}
}

func (h *Helper) SetSaveFunc(saveFunc SaveFunc) {
	h.saveFunc = saveFunc
}

func (h *Helper) HasMessages() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.messages) > 0
}

func (h *Helper) ReplaceMessages(history []StoredMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.messages = ToPromptMessages(history)
}

func (h *Helper) AddMessage(content, userName string, isUser, save bool) (*StoredMessage, error) {
	message := NewStoredMessage(h.SessionID, userName, content, isUser)

	h.mu.Lock()
	h.messages = append(h.messages, PromptMessage{
		Content: content,
		IsUser:  isUser,
	})
	h.mu.Unlock()

	if save && h.saveFunc != nil {
		if err := h.saveFunc(message); err != nil {
			return nil, err
		}
	}

	return message, nil
}

func (h *Helper) GenerateResponse(userName string, ctx context.Context, userQuestion string) (*StoredMessage, error) {
	if _, err := h.AddMessage(userQuestion, userName, true, true); err != nil {
		return nil, err
	}

	h.mu.RLock()
	messages := ToSchemaMessages(h.messages)
	h.mu.RUnlock()

	start := time.Now()
	schemaMessage, err := h.provider.GenerateResponse(ctx, messages)
	metrics.GetCollector().RecordModel(h.GetModelType(), "generate", userName, time.Since(start), err)
	if err != nil {
		return nil, err
	}

	responseMessage, err := h.AddMessage(schemaMessage.Content, userName, false, true)
	if err != nil {
		return nil, err
	}

	return responseMessage, nil
}

func (h *Helper) StreamResponse(userName string, ctx context.Context, cb StreamCallback, userQuestion string) (*StoredMessage, error) {
	if _, err := h.AddMessage(userQuestion, userName, true, true); err != nil {
		return nil, err
	}

	h.mu.RLock()
	messages := ToSchemaMessages(h.messages)
	h.mu.RUnlock()

	start := time.Now()
	content, err := h.provider.StreamResponse(ctx, messages, cb)
	metrics.GetCollector().RecordModel(h.GetModelType(), "stream", userName, time.Since(start), err)
	if err != nil {
		return nil, err
	}

	responseMessage, err := h.AddMessage(content, userName, false, true)
	if err != nil {
		return nil, err
	}

	return responseMessage, nil
}

func (h *Helper) GetModelType() string {
	return h.provider.Name()
}

type Manager struct {
	helpers map[string]map[string]*Helper
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		helpers: make(map[string]map[string]*Helper),
	}
}

func (m *Manager) GetOrCreateHelper(userName, sessionID, modelType string, config map[string]interface{}) (*Helper, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	userHelpers, exists := m.helpers[userName]
	if !exists {
		userHelpers = make(map[string]*Helper)
		m.helpers[userName] = userHelpers
	}

	helper, exists := userHelpers[sessionID]
	if exists && helper.GetModelType() == modelType {
		return helper, nil
	}

	provider, err := GetGlobalFactory().CreateProvider(context.Background(), modelType, config)
	if err != nil {
		return nil, err
	}

	helper = NewHelper(provider, sessionID)
	userHelpers[sessionID] = helper
	return helper, nil
}

func (m *Manager) RemoveHelper(userName, sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	userHelpers, exists := m.helpers[userName]
	if !exists {
		return
	}

	delete(userHelpers, sessionID)
	if len(userHelpers) == 0 {
		delete(m.helpers, userName)
	}
}

var (
	globalManager *Manager
	managerOnce   sync.Once
)

func GetGlobalManager() *Manager {
	managerOnce.Do(func() {
		globalManager = NewManager()
	})
	return globalManager
}
