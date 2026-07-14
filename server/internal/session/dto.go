package session

import "server/pkg/response"

type GetUserSessionsResponse struct {
	response.Response
	Sessions []SessionInfo `json:"sessions,omitempty"`
}

type UpdateSessionTitleRequest struct {
	SessionID string `json:"sessionId" binding:"required,max=64"`
	Title     string `json:"title" binding:"required,max=120"`
}

type UpdateSessionPinRequest struct {
	SessionID string `json:"sessionId" binding:"required,max=64"`
	Pinned    bool   `json:"pinned"`
}

type UpdateSessionArchiveRequest struct {
	SessionID string `json:"sessionId" binding:"required,max=64"`
	Archived  bool   `json:"archived"`
}
