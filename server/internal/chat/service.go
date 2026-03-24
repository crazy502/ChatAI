package chat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"server/internal/ai"
	sessionpkg "server/internal/session"
	"server/pkg/code"

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

func (s *Service) CreateSessionAndSendMessage(userName, userQuestion, modelType string) (string, string, code.Code) {
	newSession := &sessionpkg.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    sessionpkg.NormalizeTitle(userQuestion),
	}

	createdSession, err := s.sessionRepo.Create(newSession)
	if err != nil {
		log.Println("create session and send message create session error:", err)
		return "", "", code.CodeServerBusy
	}

	helper, resultCode := s.getOrCreateHydratedHelper(userName, createdSession.ID, modelType)
	if resultCode != code.CodeSuccess {
		return "", "", resultCode
	}

	aiResponse, err := helper.GenerateResponse(userName, context.Background(), userQuestion)
	if err != nil {
		log.Println("create session and send message generate response error:", err)
		return "", "", code.AIModelFail
	}

	s.touchSessionActivity(createdSession.ID)
	return createdSession.ID, aiResponse.Content, code.CodeSuccess
}

func (s *Service) CreateStreamSessionOnly(userName, userQuestion string) (string, code.Code) {
	newSession := &sessionpkg.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    sessionpkg.NormalizeTitle(userQuestion),
	}

	createdSession, err := s.sessionRepo.Create(newSession)
	if err != nil {
		log.Println("create stream session only create session error:", err)
		return "", code.CodeServerBusy
	}

	return createdSession.ID, code.CodeSuccess
}

func (s *Service) StreamMessageToExistingSession(userName, sessionID, userQuestion, modelType string, writer http.ResponseWriter) code.Code {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		log.Println("stream message unsupported")
		return code.CodeServerBusy
	}

	if _, resultCode := s.loadOwnedSession(userName, sessionID); resultCode != code.CodeSuccess {
		return resultCode
	}

	if err := writeSSEJSON(writer, flusher, map[string]bool{"ready": true}); err != nil {
		log.Println("stream message write ready error:", err)
		return code.CodeServerBusy
	}

	helper, resultCode := s.getOrCreateHydratedHelper(userName, sessionID, modelType)
	if resultCode != code.CodeSuccess {
		return resultCode
	}

	callback := func(msg string) {
		if err := writeSSEJSON(writer, flusher, map[string]string{"content": msg}); err != nil {
			log.Println("stream message write chunk error:", err)
		}
	}

	if _, err := helper.StreamResponse(userName, context.Background(), callback, userQuestion); err != nil {
		log.Println("stream response error:", err)
		return code.AIModelFail
	}

	if err := writeSSEDone(writer, flusher); err != nil {
		log.Println("stream message write done error:", err)
		return code.AIModelFail
	}

	s.touchSessionActivity(sessionID)
	return code.CodeSuccess
}

func (s *Service) ChatSend(userName, sessionID, userQuestion, modelType string) (string, code.Code) {
	if _, resultCode := s.loadOwnedSession(userName, sessionID); resultCode != code.CodeSuccess {
		return "", resultCode
	}

	helper, resultCode := s.getOrCreateHydratedHelper(userName, sessionID, modelType)
	if resultCode != code.CodeSuccess {
		return "", resultCode
	}

	aiResponse, err := helper.GenerateResponse(userName, context.Background(), userQuestion)
	if err != nil {
		log.Println("chat send generate response error:", err)
		return "", code.AIModelFail
	}

	s.touchSessionActivity(sessionID)
	return aiResponse.Content, code.CodeSuccess
}

func (s *Service) ChatStreamSend(userName, sessionID, userQuestion, modelType string, writer http.ResponseWriter) code.Code {
	return s.StreamMessageToExistingSession(userName, sessionID, userQuestion, modelType, writer)
}

func (s *Service) GetChatHistory(userName, sessionID string) ([]History, code.Code) {
	if _, resultCode := s.loadOwnedSession(userName, sessionID); resultCode != code.CodeSuccess {
		return nil, resultCode
	}

	messages, err := s.repo.GetMessagesBySessionID(sessionID)
	if err != nil {
		log.Println("get chat history error:", err)
		return nil, code.CodeServerBusy
	}

	history := make([]History, 0, len(messages))
	for _, message := range messages {
		history = append(history, History{
			IsUser:  message.IsUser,
			Content: message.Content,
		})
	}

	return history, code.CodeSuccess
}

func (s *Service) getOrCreateHydratedHelper(userName, sessionID, modelType string) (*ai.Helper, code.Code) {
	helper, err := ai.GetGlobalManager().GetOrCreateHelper(userName, sessionID, modelType, map[string]interface{}{})
	if err != nil {
		log.Println("get or create helper error:", err)
		return nil, code.AIModelFail
	}

	helper.SetSaveFunc(func(message *ai.StoredMessage) error {
		return saveWithQueue(s.repo, message)
	})

	if helper.HasMessages() {
		return helper, code.CodeSuccess
	}

	history, err := s.repo.GetMessagesBySessionID(sessionID)
	if err != nil {
		log.Println("hydrate helper load history error:", err)
		return nil, code.CodeServerBusy
	}

	if len(history) > 0 {
		helper.ReplaceMessages(toAIStoredMessages(history))
	}

	return helper, code.CodeSuccess
}

func (s *Service) loadOwnedSession(userName, sessionID string) (*sessionpkg.Session, code.Code) {
	sessionInfo, err := s.sessionRepo.GetByIDAndUserName(sessionID, userName)
	if err == gorm.ErrRecordNotFound {
		return nil, code.CodeRecordNotFound
	}
	if err != nil {
		log.Println("load owned session error:", err)
		return nil, code.CodeServerBusy
	}

	return sessionInfo, code.CodeSuccess
}

func (s *Service) touchSessionActivity(sessionID string) {
	if err := s.sessionRepo.TouchSession(sessionID, time.Now()); err != nil {
		log.Println("touch session activity error:", err)
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
