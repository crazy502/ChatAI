package session

import (
	"context"
	"strconv"

	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/requestmeta"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SessionService interface {
	ListByUserName(ctx context.Context, userName, keyword string, includeArchived bool) ([]SessionInfo, error)
	Rename(ctx context.Context, userName, sessionID, title string) error
	SetPinned(ctx context.Context, userName, sessionID string, pinned bool) error
	SetArchived(ctx context.Context, userName, sessionID string, archived bool) error
}

type Handler struct {
	service SessionService
}

func NewHandler(service SessionService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUserSessionsByUserName(c *gin.Context) {
	res := &GetUserSessionsResponse{}
	userName := c.GetString("userName")
	keyword := c.Query("keyword")
	includeArchived := parseBoolQuery(c.Query("includeArchived"))

	sessions, err := h.service.ListByUserName(c.Request.Context(), userName, keyword, includeArchived)
	if err != nil {
		response.Fail(c, err)
		return
	}

	res.Sessions = sessions
	response.OK(c, res)
}

func (h *Handler) RenameSession(c *gin.Context) {
	req := new(UpdateSessionTitleRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	ctx := attachSessionContext(c, req.SessionID)
	userName := c.GetString("userName")
	if err := h.service.Rename(ctx, userName, req.SessionID, req.Title); err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, &response.Response{})
}

func (h *Handler) UpdateSessionPin(c *gin.Context) {
	req := new(UpdateSessionPinRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	ctx := attachSessionContext(c, req.SessionID)
	userName := c.GetString("userName")
	if err := h.service.SetPinned(ctx, userName, req.SessionID, req.Pinned); err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, &response.Response{})
}

func (h *Handler) UpdateSessionArchive(c *gin.Context) {
	req := new(UpdateSessionArchiveRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	ctx := attachSessionContext(c, req.SessionID)
	userName := c.GetString("userName")
	if err := h.service.SetArchived(ctx, userName, req.SessionID, req.Archived); err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, &response.Response{})
}

func attachSessionContext(c *gin.Context, sessionID string) context.Context {
	ctx := requestmeta.WithField(c.Request.Context(), requestmeta.FieldSessionID, sessionID)
	c.Request = c.Request.WithContext(ctx)
	return ctx
}

func parseBoolQuery(value string) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return parsed
}
