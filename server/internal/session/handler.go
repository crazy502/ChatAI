package session

import (
	"net/http"
	"strconv"

	"server/pkg/code"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUserSessionsByUserName(c *gin.Context) {
	res := new(GetUserSessionsResponse)
	userName := c.GetString("userName")
	keyword := c.Query("keyword")
	includeArchived := parseBoolQuery(c.Query("includeArchived"))

	sessions, err := h.service.ListByUserName(userName, keyword, includeArchived)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	res.Success()
	res.Sessions = sessions
	c.JSON(http.StatusOK, res)
}

func (h *Handler) RenameSession(c *gin.Context) {
	req := new(UpdateSessionTitleRequest)
	res := new(response.Response)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	resultCode := h.service.Rename(userName, req.SessionID, req.Title)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateSessionPin(c *gin.Context) {
	req := new(UpdateSessionPinRequest)
	res := new(response.Response)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	resultCode := h.service.SetPinned(userName, req.SessionID, req.Pinned)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateSessionArchive(c *gin.Context) {
	req := new(UpdateSessionArchiveRequest)
	res := new(response.Response)
	userName := c.GetString("userName")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	resultCode := h.service.SetArchived(userName, req.SessionID, req.Archived)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	c.JSON(http.StatusOK, res)
}

func parseBoolQuery(value string) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return parsed
}
