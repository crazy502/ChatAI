package session

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"server/common/aihelper"
	"server/common/code"
	messageDAO "server/dao/message"
	sessionDAO "server/dao/session"
	"server/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	defaultSessionTitle  = "新会话"
	maxSessionTitleRunes = 100
)

var ctx = context.Background()

func normalizeSessionTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return defaultSessionTitle
	}

	if utf8.RuneCountInString(title) <= maxSessionTitleRunes {
		return title
	}

	runes := []rune(title)
	return strings.TrimSpace(string(runes[:maxSessionTitleRunes]))
}

func toSessionInfo(session model.Session) model.SessionInfo {
	return model.SessionInfo{
		ID:              session.ID,
		SessionID:       session.ID,
		LegacySessionID: session.ID,
		Title:           session.Title,
		LegacyTitle:     session.Title,
		Pinned:          session.Pinned,
		Archived:        session.Archived,
		LastMessageAt:   session.LastMessageAt,
		UpdatedAt:       session.UpdatedAt,
	}
}

func loadOwnedSession(userName, sessionID string) (*model.Session, code.Code) {
	sessionInfo, err := sessionDAO.GetSessionByIDAndUserName(sessionID, userName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.CodeRecordNotFound
		}
		log.Println("loadOwnedSession error:", err)
		return nil, code.CodeServerBusy
	}

	return sessionInfo, code.CodeSuccess
}

func getOrCreateHydratedHelper(userName, sessionID, modelType string) (*aihelper.AIHelper, code.Code) {
	manager := aihelper.GetGlobalManager()
	config := make(map[string]interface{})
	helper, err := manager.GetOrCreateAIHelper(userName, sessionID, modelType, config)
	if err != nil {
		log.Println("getOrCreateHydratedHelper GetOrCreateAIHelper error:", err)
		return nil, code.AIModelFail
	}

	if helper.HasMessages() {
		return helper, code.CodeSuccess
	}

	history, err := messageDAO.GetMessagesBySessionID(sessionID)
	if err != nil {
		log.Println("getOrCreateHydratedHelper GetMessagesBySessionID error:", err)
		return nil, code.CodeServerBusy
	}

	if len(history) > 0 {
		helper.ReplaceMessages(history)
	}

	return helper, code.CodeSuccess
}

func touchSessionActivity(sessionID string) {
	if err := sessionDAO.TouchSession(sessionID, time.Now()); err != nil {
		log.Println("touchSessionActivity error:", err)
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

func GetUserSessionsByUserName(userName, keyword string, includeArchived bool) ([]model.SessionInfo, error) {
	sessions, err := sessionDAO.ListSessionsByUserName(userName, keyword, includeArchived)
	if err != nil {
		return nil, err
	}

	sessionInfos := make([]model.SessionInfo, 0, len(sessions))
	for _, session := range sessions {
		sessionInfos = append(sessionInfos, toSessionInfo(session))
	}

	return sessionInfos, nil
}

func CreateSessionAndSendMessage(userName string, userQuestion string, modelType string) (string, string, code.Code) {
	newSession := &model.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    normalizeSessionTitle(userQuestion),
	}

	createdSession, err := sessionDAO.CreateSession(newSession)
	if err != nil {
		log.Println("CreateSessionAndSendMessage CreateSession error:", err)
		return "", "", code.CodeServerBusy
	}

	helper, code_ := getOrCreateHydratedHelper(userName, createdSession.ID, modelType)
	if code_ != code.CodeSuccess {
		return "", "", code_
	}

	aiResponse, err := helper.GenerateResponse(userName, ctx, userQuestion)
	if err != nil {
		log.Println("CreateSessionAndSendMessage GenerateResponse error:", err)
		return "", "", code.AIModelFail
	}

	touchSessionActivity(createdSession.ID)
	return createdSession.ID, aiResponse.Content, code.CodeSuccess
}

func CreateStreamSessionOnly(userName string, userQuestion string) (string, code.Code) {
	newSession := &model.Session{
		ID:       uuid.New().String(),
		UserName: userName,
		Title:    normalizeSessionTitle(userQuestion),
	}
	createdSession, err := sessionDAO.CreateSession(newSession)
	if err != nil {
		log.Println("CreateStreamSessionOnly CreateSession error:", err)
		return "", code.CodeServerBusy
	}
	return createdSession.ID, code.CodeSuccess
}

func StreamMessageToExistingSession(userName string, sessionID string, userQuestion string, modelType string, writer http.ResponseWriter) code.Code {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		log.Println("StreamMessageToExistingSession: streaming unsupported")
		return code.CodeServerBusy
	}

	if _, code_ := loadOwnedSession(userName, sessionID); code_ != code.CodeSuccess {
		return code_
	}

	if err := writeSSEJSON(writer, flusher, map[string]bool{"ready": true}); err != nil {
		log.Println("StreamMessageToExistingSession write ready error:", err)
		return code.CodeServerBusy
	}

	helper, code_ := getOrCreateHydratedHelper(userName, sessionID, modelType)
	if code_ != code.CodeSuccess {
		return code_
	}

	cb := func(msg string) {
		if err := writeSSEJSON(writer, flusher, map[string]string{"content": msg}); err != nil {
			log.Println("[SSE] Write error:", err)
			return
		}
	}

	if _, err := helper.StreamResponse(userName, ctx, cb, userQuestion); err != nil {
		log.Println("StreamMessageToExistingSession StreamResponse error:", err)
		return code.AIModelFail
	}

	if err := writeSSEDone(writer, flusher); err != nil {
		log.Println("StreamMessageToExistingSession write DONE error:", err)
		return code.AIModelFail
	}
	touchSessionActivity(sessionID)

	return code.CodeSuccess
}

func CreateStreamSessionAndSendMessage(userName string, userQuestion string, modelType string, writer http.ResponseWriter) (string, code.Code) {
	sessionID, code_ := CreateStreamSessionOnly(userName, userQuestion)
	if code_ != code.CodeSuccess {
		return "", code_
	}

	code_ = StreamMessageToExistingSession(userName, sessionID, userQuestion, modelType, writer)
	if code_ != code.CodeSuccess {
		return sessionID, code_
	}

	return sessionID, code.CodeSuccess
}

func ChatSend(userName string, sessionID string, userQuestion string, modelType string) (string, code.Code) {
	if _, code_ := loadOwnedSession(userName, sessionID); code_ != code.CodeSuccess {
		return "", code_
	}

	helper, code_ := getOrCreateHydratedHelper(userName, sessionID, modelType)
	if code_ != code.CodeSuccess {
		return "", code_
	}

	aiResponse, err := helper.GenerateResponse(userName, ctx, userQuestion)
	if err != nil {
		log.Println("ChatSend GenerateResponse error:", err)
		return "", code.AIModelFail
	}

	touchSessionActivity(sessionID)
	return aiResponse.Content, code.CodeSuccess
}

func GetChatHistory(userName string, sessionID string) ([]model.History, code.Code) {
	if _, code_ := loadOwnedSession(userName, sessionID); code_ != code.CodeSuccess {
		return nil, code_
	}

	messages, err := messageDAO.GetMessagesBySessionID(sessionID)
	if err != nil {
		log.Println("GetChatHistory GetMessagesBySessionID error:", err)
		return nil, code.CodeServerBusy
	}

	history := make([]model.History, 0, len(messages))
	for _, msg := range messages {
		history = append(history, model.History{
			IsUser:  msg.IsUser,
			Content: msg.Content,
		})
	}

	return history, code.CodeSuccess
}

func RenameSession(userName, sessionID, title string) code.Code {
	if _, code_ := loadOwnedSession(userName, sessionID); code_ != code.CodeSuccess {
		return code_
	}

	if err := sessionDAO.UpdateSessionTitle(sessionID, userName, normalizeSessionTitle(title)); err != nil {
		log.Println("RenameSession error:", err)
		return code.CodeServerBusy
	}

	return code.CodeSuccess
}

func SetSessionPinned(userName, sessionID string, pinned bool) code.Code {
	if _, code_ := loadOwnedSession(userName, sessionID); code_ != code.CodeSuccess {
		return code_
	}

	if err := sessionDAO.UpdateSessionPin(sessionID, userName, pinned); err != nil {
		log.Println("SetSessionPinned error:", err)
		return code.CodeServerBusy
	}

	return code.CodeSuccess
}

func SetSessionArchived(userName, sessionID string, archived bool) code.Code {
	if _, code_ := loadOwnedSession(userName, sessionID); code_ != code.CodeSuccess {
		return code_
	}

	if err := sessionDAO.UpdateSessionArchive(sessionID, userName, archived); err != nil {
		log.Println("SetSessionArchived error:", err)
		return code.CodeServerBusy
	}

	if archived {
		aihelper.GetGlobalManager().RemoveAIHelper(userName, sessionID)
	}

	return code.CodeSuccess
}

func ChatStreamSend(userName string, sessionID string, userQuestion string, modelType string, writer http.ResponseWriter) code.Code {
	return StreamMessageToExistingSession(userName, sessionID, userQuestion, modelType, writer)
}
