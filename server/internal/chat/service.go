package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"server/internal/ai"
	sessionpkg "server/internal/session"
	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/observe"
	"server/pkg/requestmeta"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo        *Repository
	sessionRepo *sessionpkg.Repository
}

func NewService(repo *Repository, sessionRepo *sessionpkg.Repository) *Service {
	return &Service{
		repo:        repo,
		sessionRepo: sessionRepo,
	}
}

func (s *Service) CreateSessionAndSendMessage(ctx context.Context, userName, userQuestion, modelType string) (string, string, error) {
	//1. 创建会话
	createdSession, ctx, err := s.createSession(ctx, userName, userQuestion)
	if err != nil {
		return "", "", err
	}

	//2. 获取或创建助手实例
	helper, err := s.getOrCreateHydratedHelper(ctx, userName, createdSession.ID, modelType)
	if err != nil {
		return "", "", err
	}

	//3. 生成AI响应
	aiResponse, err := helper.GenerateResponse(userName, ctx, userQuestion)
	if err != nil {
		return "", "", apperror.Wrap(code.AIModelFail, err, "generate chat response failed").
			WithField("session_id", createdSession.ID).
			WithField("model_type", modelType)
	}

	//4. 触发会话活动更新
	s.touchSessionActivity(ctx, createdSession.ID)
	return createdSession.ID, aiResponse.Content, nil
}

func (s *Service) CreateStreamSessionOnly(ctx context.Context, userName, userQuestion string) (string, error) {
	createdSession, _, err := s.createSession(ctx, userName, userQuestion)
	if err != nil {
		return "", err
	}

	return createdSession.ID, nil
}

func (s *Service) StreamMessageToExistingSession(ctx context.Context, userName, sessionID, userQuestion, modelType string, writer http.ResponseWriter) error {
	//1. 检查是否支持流式响应
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return apperror.New(code.CodeServerBusy, "streaming response is not supported").
			WithField("session_id", sessionID)
	}

	//2. 加载会话
	ctx = requestmeta.WithField(ctx, requestmeta.FieldSessionID, sessionID)
	//3. 加载助手实例
	if _, err := s.loadOwnedSession(ctx, userName, sessionID); err != nil {
		return err
	}

	//4. 写入流式响应事件：ready
	if err := writeSSEJSON(writer, flusher, map[string]bool{"ready": true}); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "write stream ready event failed").
			WithField("session_id", sessionID)
	}

	helper, err := s.getOrCreateHydratedHelper(ctx, userName, sessionID, modelType)
	if err != nil {
		return err
	}

	callback := func(msg string) {
		if writeErr := writeSSEJSON(writer, flusher, map[string]string{"content": msg}); writeErr != nil {
			observe.Error(ctx, "write stream chunk failed", apperror.Wrap(code.CodeServerBusy, writeErr, "write stream chunk failed"))
		}
	}

	if _, err := helper.StreamResponse(userName, ctx, callback, userQuestion); err != nil {
		return apperror.Wrap(code.AIModelFail, err, "stream chat response failed").
			WithField("session_id", sessionID).
			WithField("model_type", modelType)
	}

	if err := writeSSEDone(writer, flusher); err != nil {
		return apperror.Wrap(code.AIModelFail, err, "write stream done event failed").
			WithField("session_id", sessionID)
	}

	s.touchSessionActivity(ctx, sessionID)
	return nil
}

// ChatSend 发送聊天消息
// return: 响应内容
// err: 错误
func (s *Service) ChatSend(ctx context.Context, userName, sessionID, userQuestion, modelType string) (string, error) {
	//1. 加载会话
	ctx = requestmeta.WithField(ctx, requestmeta.FieldSessionID, sessionID)

	//2. 加载助手实例
	if _, err := s.loadOwnedSession(ctx, userName, sessionID); err != nil {
		return "", err
	}

	//3. 获取或创建助手实例
	helper, err := s.getOrCreateHydratedHelper(ctx, userName, sessionID, modelType)
	if err != nil {
		return "", err
	}

	//4. 生成AI响应
	aiResponse, err := helper.GenerateResponse(userName, ctx, userQuestion)
	if err != nil {
		return "", apperror.Wrap(code.AIModelFail, err, "generate chat response failed").
			WithField("session_id", sessionID).
			WithField("model_type", modelType)
	}

	s.touchSessionActivity(ctx, sessionID)
	return aiResponse.Content, nil
}

func (s *Service) ChatStreamSend(ctx context.Context, userName, sessionID, userQuestion, modelType string, writer http.ResponseWriter) error {
	return s.StreamMessageToExistingSession(ctx, userName, sessionID, userQuestion, modelType, writer)
}

func (s *Service) GetChatHistory(ctx context.Context, userName, sessionID string) ([]History, error) {
	ctx = requestmeta.WithField(ctx, requestmeta.FieldSessionID, sessionID)

	if _, err := s.loadOwnedSession(ctx, userName, sessionID); err != nil {
		return nil, err
	}

	messages, err := s.repo.GetMessagesBySessionID(sessionID)
	if err != nil {
		return nil, apperror.Wrap(code.CodeServerBusy, err, "load chat history failed").
			WithField("session_id", sessionID)
	}

	history := make([]History, 0, len(messages))
	for _, message := range messages {
		history = append(history, History{
			IsUser:  message.IsUser,
			Content: message.Content,
		})
	}

	return history, nil
}

func (s *Service) createSession(ctx context.Context, userName, userQuestion string) (*sessionpkg.Session, context.Context, error) {
	newSession := &sessionpkg.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    sessionpkg.NormalizeTitle(userQuestion),
	}

	createdSession, err := s.sessionRepo.Create(newSession)
	if err != nil {
		return nil, ctx, apperror.Wrap(code.CodeServerBusy, err, "create session failed").
			WithField("user_name", userName)
	}

	ctx = requestmeta.WithField(ctx, requestmeta.FieldSessionID, createdSession.ID)
	return createdSession, ctx, nil
}

func (s *Service) getOrCreateHydratedHelper(ctx context.Context, userName, sessionID, modelType string) (*ai.Helper, error) {
	//1. 创建或获取助手实例
	helper, err := ai.GetGlobalManager().GetOrCreateHelper(userName, sessionID, modelType, map[string]interface{}{})
	if err != nil {
		return nil, mapModelFactoryError(err).
			WithField("session_id", sessionID).
			WithField("model_type", modelType)
	}

	//2. 设置消息保存函数
	helper.SetSaveFunc(func(message *ai.StoredMessage) error {
		return saveWithQueue(s.repo, message)
	})

	//3. 加载历史消息（如果需要）
	if helper.HasMessages() {
		return helper, nil
	}

	history, err := s.repo.GetMessagesBySessionID(sessionID)
	if err != nil {
		return nil, apperror.Wrap(code.CodeServerBusy, err, "hydrate helper history failed").
			WithField("session_id", sessionID)
	}

	if len(history) > 0 {
		helper.ReplaceMessages(toAIStoredMessages(history))
	}

	return helper, nil
}

func (s *Service) loadOwnedSession(ctx context.Context, userName, sessionID string) (*sessionpkg.Session, error) {
	sessionInfo, err := s.sessionRepo.GetByIDAndUserName(sessionID, userName)
	if err == gorm.ErrRecordNotFound {
		return nil, apperror.New(code.CodeRecordNotFound, code.CodeRecordNotFound.Msg()).
			WithField("user_name", userName).
			WithField("session_id", sessionID)
	}
	if err != nil {
		return nil, apperror.Wrap(code.CodeServerBusy, err, "load session failed").
			WithField("user_name", userName).
			WithField("session_id", sessionID)
	}

	return sessionInfo, nil
}

// touchSessionActivity 触发会话活动更新
// return: 错误
func (s *Service) touchSessionActivity(ctx context.Context, sessionID string) {
	if err := s.sessionRepo.TouchSession(sessionID, time.Now()); err != nil {
		observe.Error(ctx, "touch session activity failed", apperror.Wrap(code.CodeServerBusy, err, "touch session activity failed").
			WithField("session_id", sessionID))
	}
}

func mapModelFactoryError(err error) *apperror.Error {
	if err == nil {
		return nil
	}

	message := err.Error()
	switch {
	case strings.Contains(message, "unsupported model type"):
		return apperror.Wrap(code.AIModelNotFind, err, code.AIModelNotFind.Msg())
	default:
		return apperror.Wrap(code.AIModelCannotOpen, err, code.AIModelCannotOpen.Msg())
	}
}

func writeSSEJSON(writer http.ResponseWriter, flusher http.Flusher, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err := writer.Write([]byte("data: " + string(data) + "\n\n")); err != nil {
		return err
	}

	flusher.Flush()
	return nil
}

func writeSSEDone(writer http.ResponseWriter, flusher http.Flusher) error {
	if _, err := writer.Write([]byte("data: [DONE]\n\n")); err != nil {
		return err
	}

	flusher.Flush()
	return nil
}
