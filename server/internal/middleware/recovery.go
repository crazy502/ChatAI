package middleware

import (
	"server/pkg/apperror"
	"server/pkg/observe"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				appErr := apperror.Panic(recovered)

				if c.Writer.Written() && isStreamingPath(c.Request.URL.Path) {
					observe.Error(c.Request.Context(), "panic recovered during streaming response", appErr)
					response.SSEError(c, appErr)
					c.Abort()
					return
				}

				if c.Writer.Written() {
					observe.Error(c.Request.Context(), "panic recovered after response write", appErr)
					c.Abort()
					return
				}

				response.Fail(c, appErr)
				c.Abort()
			}
		}()

		c.Next()
	}
}
