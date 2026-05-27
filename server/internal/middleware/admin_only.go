package middleware

import (
	"strings"

	"server/infra/config"
	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireAdmin 管理员权限验证中间件
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. 从Gin上下文中获取 isAdmin 标记和用户名
		adminUsername := strings.TrimSpace(config.GetConfig().AdminConfig.Username)
		//2. 检查管理员用户名是否为空
		if adminUsername == "" {
			adminUsername = "admin@qq.com"
		}
		if !c.GetBool("isAdmin") || c.GetString("userName") != adminUsername {
			response.Fail(c, apperror.New(code.CodeForbidden, code.CodeForbidden.Msg()))
			c.Abort()
			return
		}

		c.Next()
	}
}
