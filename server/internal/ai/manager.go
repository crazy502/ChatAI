package ai

import (
	"context"
	"sync"
	"time"

	"server/infra/metrics"

	"github.com/cloudwego/eino/schema"
)

type SaveFunc func(*StoredMessage) error

type LoadHistoryFunc func() ([]StoredMessage, error)

const maxContextMessages = 100

type Helper struct {
	turnMu    sync.Mutex
	hydrateMu sync.Mutex
	hydrated  bool
	provider  Provider        // 模型提供方
	messages  []PromptMessage // 消息队列
	mu        sync.RWMutex    // 读写锁，用于保护消息队列
	SessionID string          // 会话ID
	saveFunc  SaveFunc        // 保存消息函数
}

func NewHelper(provider Provider, sessionID string) *Helper {
	return &Helper{
		provider:  provider,
		messages:  make([]PromptMessage, 0),
		SessionID: sessionID,
	}
}

func (h *Helper) SetSaveFunc(saveFunc SaveFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.saveFunc = saveFunc
}

func (h *Helper) HasMessages() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.messages) > 0
}

// ReplaceMessages 替换助手队列中的消息
// history: 历史消息队列
func (h *Helper) ReplaceMessages(history []StoredMessage) {
	//1. 加锁并替换消息队列
	h.mu.Lock()
	//2. 解锁
	defer h.mu.Unlock()
	//3. 替换消息队列
	h.messages = ToPromptMessages(history)
	h.hydrated = true
}

func (h *Helper) EnsureHydrated(load LoadHistoryFunc) error {
	h.hydrateMu.Lock()
	defer h.hydrateMu.Unlock()

	h.mu.RLock()
	alreadyHydrated := h.hydrated
	h.mu.RUnlock()
	if alreadyHydrated {
		return nil
	}

	history, err := load()
	if err != nil {
		return err
	}

	h.mu.Lock()
	h.messages = ToPromptMessages(history)
	h.hydrated = true
	h.mu.Unlock()
	return nil
}

// AddMessage 添加消息到助手队列
// return: 消息实例
// err: 错误
func (h *Helper) AddMessage(content, userName string, isUser, save bool) (*StoredMessage, error) {
	//1. 创建消息实例
	message := NewStoredMessage(h.SessionID, userName, content, isUser)

	// Persist first so failed writes never leak into the model context.
	h.mu.RLock()
	saveFunc := h.saveFunc
	h.mu.RUnlock()
	if save && saveFunc != nil {
		if err := saveFunc(message); err != nil {
			return nil, err
		}
	}

	//2. 加锁并添加消息到队列
	h.mu.Lock()
	h.messages = append(h.messages, PromptMessage{
		Content: content,
		IsUser:  isUser,
	})
	h.hydrated = true
	//3. 解锁
	h.mu.Unlock()

	//4. 保存消息
	return message, nil
}

func (h *Helper) GenerateResponse(userName string, ctx context.Context, userQuestion string) (*StoredMessage, error) {
	h.turnMu.Lock()
	defer h.turnMu.Unlock()

	if _, err := h.AddMessage(userQuestion, userName, true, true); err != nil {
		return nil, err
	}

	messages := h.contextMessages()

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
	h.turnMu.Lock()
	defer h.turnMu.Unlock()

	if _, err := h.AddMessage(userQuestion, userName, true, true); err != nil {
		return nil, err
	}

	messages := h.contextMessages()

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

func (h *Helper) contextMessages() []*schema.Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	start := 0
	if len(h.messages) > maxContextMessages {
		start = len(h.messages) - maxContextMessages
	}
	return ToSchemaMessages(h.messages[start:])
}

type Manager struct {
	helpers map[string]map[string]*Helper // 用户会话ID到助手的映射
	mu      sync.RWMutex                  // 读写锁，用于保护助手映射
}

func NewManager() *Manager {
	return &Manager{
		helpers: make(map[string]map[string]*Helper),
	}
}

func (m *Manager) GetOrCreateHelper(userName, sessionID, modelType string, config map[string]interface{}) (*Helper, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//1. 检查用户的Helper映射是否存在
	userHelpers, exists := m.helpers[userName]
	if !exists {
		userHelpers = make(map[string]*Helper)
		m.helpers[userName] = userHelpers
	}

	//2. 检查会话的Helper实例是否存在且模型类型匹配
	helper, exists := userHelpers[sessionID]
	if exists && helper.GetModelType() == modelType {
		return helper, nil
	}

	//3. 创建新的模型提供者实例
	provider, err := GetGlobalFactory().CreateProvider(context.Background(), modelType, config)
	if err != nil {
		return nil, err
	}

	//4. 创建并缓存新的Helper实例
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

// GetGlobalManager 获取全局助手管理器
// return: 全局助手管理器实例
func GetGlobalManager() *Manager {
	managerOnce.Do(func() {
		globalManager = NewManager()
	})
	return globalManager
}
