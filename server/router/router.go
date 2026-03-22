package router

import (
	"server/middleware/jwt"
	"server/middleware/observability"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()
	r.Use(observability.RequestMetrics())
	enterRouter := r.Group("/api/v1")
	{
		RegisterUserRouter(enterRouter.Group("/user"))
	}
	//后续登录的接口需要jwt鉴权
	{
		AIGroup := enterRouter.Group("/AI")
		AIGroup.Use(jwt.Auth())
		AIRouter(AIGroup)
	}

	{
		adminGroup := enterRouter.Group("/admin")
		adminGroup.Use(jwt.Auth())
		adminGroup.Use(jwt.RequireAdmin())
		RegisterAdminRouter(adminGroup)
	}

	return r
}
