package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AllMetrics(c *gin.Context) {
	res := new(AllMetricsResponse)
	res.Success()
	res.Snapshot = h.service.AllMetricsSnapshot()
	c.JSON(http.StatusOK, res)
}
