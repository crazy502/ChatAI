package chat

import "server/pkg/response"

type CreateSessionAndSendMessageRequest struct {
	UserQuestion string `json:"question" binding:"required,max=12000"`
	ModelType    string `json:"modelType" binding:"required,oneof=qwen deepseek"`
}

type CreateSessionAndSendMessageResponse struct {
	AiInformation string `json:"Information,omitempty"`
	SessionID     string `json:"sessionId,omitempty"`
	response.Response
}

type ChatSendRequest struct {
	UserQuestion string `json:"question" binding:"required,max=12000"`
	ModelType    string `json:"modelType" binding:"required,oneof=qwen deepseek"`
	SessionID    string `json:"sessionId,omitempty" binding:"required,max=64"`
}

type ChatSendResponse struct {
	AiInformation string `json:"Information,omitempty"`
	response.Response
}

type ChatHistoryRequest struct {
	SessionID string `json:"sessionId,omitempty" binding:"required,max=64"`
}

type ChatHistoryResponse struct {
	History []History `json:"history"`
	response.Response
}
