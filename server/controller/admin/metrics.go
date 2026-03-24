package admin

import (
	"net/http"

	"server/common/metrics"
	"server/controller"

	"github.com/gin-gonic/gin"
)

type AllMetricsResponse struct {
	controller.Response
	Snapshot metrics.AllMetricsSnapshot `json:"snapshot"`
}

func AllMetrics(c *gin.Context) {
	res := new(AllMetricsResponse)
	res.Success()
	res.Snapshot = metrics.GetCollector().AllMetricsSnapshot()
	c.JSON(http.StatusOK, res)
}
