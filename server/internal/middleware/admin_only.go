package middleware

import (
	"strings"

	"server/infra/config"
	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminUsername := strings.TrimSpace(config.GetConfig().AdminConfig.Username)
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
