package middleware

import (
	"net/http"
	"strconv"
	"time"

	"server/infra/cache"
	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

func MaxBodyBytes(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil && limit > 0 {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		}
		c.Next()
	}
}

func RateLimit(namespace string, limit int64, window time.Duration, authenticated bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.ClientIP()
		if authenticated {
			identity = c.GetString("userName")
		}

		allowed, retryAfter, err := cache.AllowRequest(c.Request.Context(), namespace, identity, limit, window)
		if err != nil {
			response.Fail(c, apperror.Wrap(code.CodeServerBusy, err, "check request rate limit failed"))
			c.Abort()
			return
		}
		if !allowed {
			seconds := int64(retryAfter.Round(time.Second) / time.Second)
			if seconds < 1 {
				seconds = 1
			}
			c.Header("Retry-After", strconv.FormatInt(seconds, 10))
			response.Fail(c, apperror.New(code.CodeTooManyRequests, code.CodeTooManyRequests.Msg()))
			c.Abort()
			return
		}

		c.Next()
	}
}
