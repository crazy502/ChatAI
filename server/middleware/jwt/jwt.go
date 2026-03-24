package jwt

import (
	"net/http"
	"strings"

	"server/common/code"
	"server/config"
	"server/controller"
	"server/utils/myjwt"

	"github.com/gin-gonic/gin"
)

// 璇诲彇jwt
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(controller.Response)

		authHeader := c.GetHeader("Authorization")
		token := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		}
		if token == "" {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		claims, ok := myjwt.ParseToken(token)
		if !ok {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)
		c.Set("userName", claims.Username)
		c.Set("isAdmin", claims.IsAdmin)
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(controller.Response)
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
