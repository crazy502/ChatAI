package middleware

import (
	"net/http"
	"strings"

	"server/infra/config"
	"server/pkg/code"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(response.Response)
		adminUsername := strings.TrimSpace(config.GetConfig().AdminConfig.Username)
		if adminUsername == "" {
			adminUsername = "admin"
		}

		if !c.GetBool("isAdmin") || c.GetString("userName") != adminUsername {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}
