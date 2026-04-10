package admin

import (
	"context"

	"server/infra/metrics"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminService interface {
	AllMetricsSnapshot(ctx context.Context) metrics.AllMetricsSnapshot
}

type Handler struct {
	service AdminService
}

func NewHandler(service AdminService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AllMetrics(c *gin.Context) {
	res := &AllMetricsResponse{
		Snapshot: h.service.AllMetricsSnapshot(c.Request.Context()),
	}
	response.OK(c, res)
}
