package middleware

import (
	"strings"

	"server/pkg/apperror"
	"server/pkg/code"
	tokenjwt "server/pkg/jwt"
	"server/pkg/requestmeta"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		}
		if token == "" {
			response.Fail(c, apperror.New(code.CodeInvalidToken, code.CodeInvalidToken.Msg()))
			c.Abort()
			return
		}

		claims, ok := tokenjwt.ParseToken(token)
		if !ok {
			response.Fail(c, apperror.New(code.CodeInvalidToken, code.CodeInvalidToken.Msg()))
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)
		c.Set("userName", claims.Username)
		c.Set("isAdmin", claims.IsAdmin)

		ctx := requestmeta.WithFields(c.Request.Context(), requestmeta.Fields{
			requestmeta.FieldUserID:   claims.ID,
			requestmeta.FieldUserName: claims.Username,
			requestmeta.FieldIsAdmin:  claims.IsAdmin,
		})
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
