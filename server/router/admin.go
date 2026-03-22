package router

import (
	admin "server/controller/admin"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.RouterGroup) {
	{
		r.GET("/metrics/overview", admin.MetricsOverview)
		r.GET("/metrics/routes", admin.RouteMetrics)
		r.GET("/metrics/users", admin.UserMetrics)
		r.GET("/metrics/business-codes", admin.BusinessCodeMetrics)
		r.GET("/metrics/models", admin.ModelMetrics)
		r.GET("/metrics/model-failures", admin.ModelFailureMetrics)
	}
}
