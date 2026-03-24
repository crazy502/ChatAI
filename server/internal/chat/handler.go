package chat

import (
	"fmt"
	"net/http"

	"server/pkg/code"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateSessionAndSendMessage(c *gin.Context) {
	req := new(CreateSessionAndSendMessageRequest)
	res := new(CreateSessionAndSendMessageResponse)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	sessionID, aiInformation, resultCode := h.service.CreateSessionAndSendMessage(userName, req.UserQuestion, req.ModelType)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	res.AiInformation = aiInformation
	res.SessionID = sessionID
	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateStreamSessionAndSendMessage(c *gin.Context) {
	req := new(CreateSessionAndSendMessageRequest)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid parameters"})
		return
	}

	h.prepareSSE(c)

	sessionID, resultCode := h.service.CreateStreamSessionOnly(userName, req.UserQuestion)
	if resultCode != code.CodeSuccess {
		c.SSEvent("error", gin.H{"message": "Failed to create session"})
		return
	}

	_, _ = c.Writer.WriteString(fmt.Sprintf("data: {\"sessionId\": \"%s\"}\n\n", sessionID))
	c.Writer.Flush()

	resultCode = h.service.StreamMessageToExistingSession(userName, sessionID, req.UserQuestion, req.ModelType, http.ResponseWriter(c.Writer))
	if resultCode != code.CodeSuccess {
		c.SSEvent("error", gin.H{"message": "Failed to send message"})
		return
	}
}

func (h *Handler) ChatSend(c *gin.Context) {
	req := new(ChatSendRequest)
	res := new(ChatSendResponse)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	aiInformation, resultCode := h.service.ChatSend(userName, req.SessionID, req.UserQuestion, req.ModelType)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	res.AiInformation = aiInformation
	c.JSON(http.StatusOK, res)
}

func (h *Handler) ChatStreamSend(c *gin.Context) {
	req := new(ChatSendRequest)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid parameters"})
		return
	}

	h.prepareSSE(c)

	resultCode := h.service.ChatStreamSend(userName, req.SessionID, req.UserQuestion, req.ModelType, http.ResponseWriter(c.Writer))
	if resultCode != code.CodeSuccess {
		c.SSEvent("error", gin.H{"message": "Failed to send message"})
		return
	}
}

func (h *Handler) ChatHistory(c *gin.Context) {
	req := new(ChatHistoryRequest)
	res := new(ChatHistoryResponse)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	history, resultCode := h.service.GetChatHistory(userName, req.SessionID)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	res.History = history
	c.JSON(http.StatusOK, res)
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
