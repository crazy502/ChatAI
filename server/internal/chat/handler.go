package chat

import (
	"context"
	"fmt"
	"net/http"

	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/requestmeta"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChatService interface {
	CreateSessionAndSendMessage(ctx context.Context, userName, userQuestion, modelType string) (string, string, error)
	CreateStreamSessionOnly(ctx context.Context, userName, userQuestion string) (string, error)
	StreamMessageToExistingSession(ctx context.Context, userName, sessionID, userQuestion, modelType string, writer http.ResponseWriter) error
	ChatSend(ctx context.Context, userName, sessionID, userQuestion, modelType string) (string, error)
	ChatStreamSend(ctx context.Context, userName, sessionID, userQuestion, modelType string, writer http.ResponseWriter) error
	GetChatHistory(ctx context.Context, userName, sessionID string) ([]History, error)
}

type Handler struct {
	service ChatService
}

func NewHandler(service ChatService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateSessionAndSendMessage(c *gin.Context) {
	req := new(CreateSessionAndSendMessageRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	userName := c.GetString("userName")
	sessionID, aiInformation, err := h.service.CreateSessionAndSendMessage(c.Request.Context(), userName, req.UserQuestion, req.ModelType)
	if err != nil {
		response.Fail(c, err)
		return
	}

	c.Request = c.Request.WithContext(requestmeta.WithField(c.Request.Context(), requestmeta.FieldSessionID, sessionID))
	res := &CreateSessionAndSendMessageResponse{
		AiInformation: aiInformation,
		SessionID:     sessionID,
	}
	response.OK(c, res)
}

func (h *Handler) CreateStreamSessionAndSendMessage(c *gin.Context) {
	req := new(CreateSessionAndSendMessageRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	h.prepareSSE(c)
	userName := c.GetString("userName")

	sessionID, err := h.service.CreateStreamSessionOnly(c.Request.Context(), userName, req.UserQuestion)
	if err != nil {
		response.SSEError(c, err)
		return
	}

	c.Request = c.Request.WithContext(requestmeta.WithField(c.Request.Context(), requestmeta.FieldSessionID, sessionID))

	_, _ = c.Writer.WriteString(fmt.Sprintf("data: {\"sessionId\": \"%s\"}\n\n", sessionID))
	c.Writer.Flush()

	if err := h.service.StreamMessageToExistingSession(c.Request.Context(), userName, sessionID, req.UserQuestion, req.ModelType, http.ResponseWriter(c.Writer)); err != nil {
		response.SSEError(c, err)
		return
	}
}

func (h *Handler) ChatSend(c *gin.Context) {
	req := new(ChatSendRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	ctx := attachSessionContext(c, req.SessionID)
	userName := c.GetString("userName")

	aiInformation, err := h.service.ChatSend(ctx, userName, req.SessionID, req.UserQuestion, req.ModelType)
	if err != nil {
		response.Fail(c, err)
		return
	}

	res := &ChatSendResponse{
		AiInformation: aiInformation,
	}
	response.OK(c, res)
}

func (h *Handler) ChatStreamSend(c *gin.Context) {
	req := new(ChatSendRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	ctx := attachSessionContext(c, req.SessionID)
	userName := c.GetString("userName")
	h.prepareSSE(c)

	if err := h.service.ChatStreamSend(ctx, userName, req.SessionID, req.UserQuestion, req.ModelType, http.ResponseWriter(c.Writer)); err != nil {
		response.SSEError(c, err)
		return
	}
}

func (h *Handler) ChatHistory(c *gin.Context) {
	req := new(ChatHistoryRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	ctx := attachSessionContext(c, req.SessionID)
	userName := c.GetString("userName")
	history, err := h.service.GetChatHistory(ctx, userName, req.SessionID)
	if err != nil {
		response.Fail(c, err)
		return
	}

	res := &ChatHistoryResponse{
		History: history,
	}
	response.OK(c, res)
}

func (h *Handler) prepareSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache, no-transform")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	c.Writer.WriteHeaderNow()
}

func attachSessionContext(c *gin.Context, sessionID string) context.Context {
	ctx := requestmeta.WithField(c.Request.Context(), requestmeta.FieldSessionID, sessionID)
	c.Request = c.Request.WithContext(ctx)
	return ctx
}
