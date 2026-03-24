package router

import (
	admin "server/controller/admin"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.RouterGroup) {
	{
		r.GET("/metrics/all", admin.AllMetrics)
	}
}
