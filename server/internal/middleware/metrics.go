package middleware

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"server/infra/metrics"
	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/observe"

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

func RequestObserver() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		shouldCapture := !isStreamingPath(path)

		var writer *bodyCaptureWriter
		if shouldCapture {
			writer = &bodyCaptureWriter{
				ResponseWriter: c.Writer,
				limit:          maxCaptureBytes,
			}
			c.Writer = writer
		}

		c.Next()

		routePath := c.FullPath()
		if routePath == "" {
			routePath = path
		}

		latency := time.Since(start)
		businessCode := int64(0)
		if writer != nil && isJSONResponse(c.Writer.Header().Get("Content-Type")) {
			businessCode = parseBusinessCode(writer.body.Bytes())
		}

		logRequest(c, routePath, latency, businessCode)

		if strings.HasPrefix(path, "/api/v1/admin/metrics") {
			return
		}

		metrics.GetCollector().RecordRequest(
			c.Request.Method,
			routePath,
			c.GetString("userName"),
			latency,
			businessCode,
			c.Writer.Status(),
		)
	}
}

func isStreamingPath(path string) bool {
	return strings.HasSuffix(path, "/chat/send-stream") ||
		strings.HasSuffix(path, "/chat/send-stream-new-session")
}

func isJSONResponse(contentType string) bool {
	return strings.Contains(contentType, "application/json")
}

func logRequest(c *gin.Context, routePath string, latency time.Duration, businessCode int64) {
	attrs := []any{
		"method", c.Request.Method,
		"path", routePath,
		"http_status", c.Writer.Status(),
		"latency_ms", latency.Milliseconds(),
	}
	if businessCode != 0 {
		attrs = append(attrs, "business_code", businessCode)
	}

	if c.Writer.Status() >= 500 || businessCode >= code.CodeServerBusy.Code() {
		serverErr := apperror.New(code.CodeServerBusy, code.CodeServerBusy.Msg())
		if businessCode >= code.CodeServerBusy.Code() {
			serverErr = apperror.New(code.Code(businessCode), code.Code(businessCode).Msg())
		}
		observe.Error(c.Request.Context(), "request completed with server error", serverErr, attrs...)
		return
	}

	if c.Writer.Status() >= 400 || (businessCode != 0 && businessCode != code.CodeSuccess.Code()) {
		observe.Warn(c.Request.Context(), "request completed with business error", attrs...)
		return
	}

	observe.Info(c.Request.Context(), "request completed", attrs...)
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
