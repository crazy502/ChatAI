package jwt

import (
	"net/http"
	"strings"

	"server/common/code"
	"server/common/mysql"
	"server/controller"
	"server/utils/myjwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeForbidden))
			c.Abort()
			return
		}

		userIDValue, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeForbidden))
			c.Abort()
			return
		}

		currentUser, err := mysql.GetUserByID(userIDValue)
		if err == gorm.ErrRecordNotFound || currentUser == nil || !currentUser.IsAdmin {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeForbidden))
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
			c.Abort()
			return
		}

		c.Set("isAdmin", true)
		c.Next()
	}
}
