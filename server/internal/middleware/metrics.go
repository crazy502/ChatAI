package middleware

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"server/infra/metrics"

	"github.com/gin-gonic/gin"
)

const maxCaptureBytes = 4096

type bodyCaptureWriter struct {
	gin.ResponseWriter
	body  bytes.Buffer
	limit int
}

func (w *bodyCaptureWriter) Write(data []byte) (int, error) {
	w.captureBytes(data)
	return w.ResponseWriter.Write(data)
}

func (w *bodyCaptureWriter) WriteString(s string) (int, error) {
	w.captureBytes([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

func (w *bodyCaptureWriter) captureBytes(data []byte) {
	if w.limit <= 0 || len(data) == 0 {
		return
	}

	remaining := w.limit - w.body.Len()
	if remaining <= 0 {
		return
	}

	if len(data) > remaining {
		data = data[:remaining]
	}

	_, _ = w.body.Write(data)
}

func RequestMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/v1/admin/metrics") || isStreamingPath(path) {
			c.Next()
			return
		}

		writer := &bodyCaptureWriter{
			ResponseWriter: c.Writer,
			limit:          maxCaptureBytes,
		}
		c.Writer = writer

		c.Next()

		routePath := c.FullPath()
		if routePath == "" {
			routePath = path
		}

		businessCode := int64(0)
		contentType := c.Writer.Header().Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			businessCode = parseBusinessCode(writer.body.Bytes())
		}

		metrics.GetCollector().RecordRequest(
			c.Request.Method,
			routePath,
			c.GetString("userName"),
			time.Since(start),
			businessCode,
			c.Writer.Status(),
		)
	}
}

func isStreamingPath(path string) bool {
	return strings.HasSuffix(path, "/chat/send-stream") ||
		strings.HasSuffix(path, "/chat/send-stream-new-session")
}

func parseBusinessCode(body []byte) int64 {
	if len(body) == 0 {
		return 0
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0
	}

	rawCode, exists := payload["status_code"]
	if !exists {
		return 0
	}

	var codeValue int64
	if err := json.Unmarshal(rawCode, &codeValue); err == nil {
		return codeValue
	}

	var stringCode string
	if err := json.Unmarshal(rawCode, &stringCode); err != nil {
		return 0
	}

	parsed, err := strconv.ParseInt(stringCode, 10, 64)
	if err != nil {
		return 0
	}

	return parsed
}
