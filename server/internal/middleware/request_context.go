package middleware

import (
	"strings"

	"server/pkg/requestmeta"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := requestmeta.WithField(c.Request.Context(), requestmeta.FieldRequestID, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
