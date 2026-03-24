package middleware

import (
	"net/http"
	"strings"

	"server/pkg/code"
	tokenjwt "server/pkg/jwt"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(response.Response)

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

		claims, ok := tokenjwt.ParseToken(token)
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
