package session

import "server/pkg/response"

type GetUserSessionsResponse struct {
	response.Response
	Sessions []SessionInfo `json:"sessions,omitempty"`
}

type UpdateSessionTitleRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Title     string `json:"title" binding:"required"`
}

type UpdateSessionPinRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Pinned    bool   `json:"pinned"`
}

type UpdateSessionArchiveRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Archived  bool   `json:"archived"`
}
