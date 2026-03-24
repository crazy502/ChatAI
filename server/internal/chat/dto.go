package chat

import "server/pkg/response"

type CreateSessionAndSendMessageRequest struct {
	UserQuestion string `json:"question" binding:"required"`
	ModelType    string `json:"modelType" binding:"required"`
}

type CreateSessionAndSendMessageResponse struct {
	AiInformation string `json:"Information,omitempty"`
	SessionID     string `json:"sessionId,omitempty"`
	response.Response
}

type ChatSendRequest struct {
	UserQuestion string `json:"question" binding:"required"`
	ModelType    string `json:"modelType" binding:"required"`
	SessionID    string `json:"sessionId,omitempty" binding:"required"`
}

type ChatSendResponse struct {
	AiInformation string `json:"Information,omitempty"`
	response.Response
}

type ChatHistoryRequest struct {
	SessionID string `json:"sessionId,omitempty" binding:"required"`
}

type ChatHistoryResponse struct {
	History []History `json:"history"`
	response.Response
}
