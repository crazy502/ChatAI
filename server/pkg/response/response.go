package response

import (
	"fmt"
	"net/http"

	"server/infra/metrics"
	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/observe"
	"server/pkg/requestmeta"

	"github.com/gin-gonic/gin"
)

type ErrorDetail struct {
	Code    int64             `json:"code"`
	Message string            `json:"message"`
	Stack   []apperror.Frame  `json:"stack,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type Response struct {
	StatusCode code.Code    `json:"status_code"`
	StatusMsg  string       `json:"status_msg,omitempty"`
	RequestID  string       `json:"request_id,omitempty"`
	Error      *ErrorDetail `json:"error,omitempty"`
}

type envelope interface {
	Base() *Response
}

func (r *Response) Base() *Response {
	if r == nil {
		return nil
	}
	return r
}

func (r *Response) CodeOf(c code.Code) Response {
	if r == nil {
		r = new(Response)
	}

	r.StatusCode = c
	r.StatusMsg = c.Msg()
	r.Error = nil
	return *r
}

func (r *Response) Success() {
	if r == nil {
		return
	}

	r.StatusCode = code.CodeSuccess
	r.StatusMsg = code.CodeSuccess.Msg()
	r.Error = nil
}

func OK(c *gin.Context, payload envelope) {
	if payload == nil {
		return
	}

	base := payload.Base()
	if base != nil {
		base.Success()
		base.RequestID = RequestID(c)
	}

	c.JSON(http.StatusOK, payload)
}

func Fail(c *gin.Context, err error) {
	appErr := apperror.From(err)
	metrics.GetCollector().RecordError(c.Request.Context(), appErr)

	logFailure(c, appErr)

	res := Response{
		StatusCode: appErr.Code,
		StatusMsg:  appErr.Code.Msg(),
		RequestID:  RequestID(c),
		Error:      buildErrorDetail(appErr),
	}
	c.JSON(http.StatusOK, res)
}

func SSEError(c *gin.Context, err error) {
	appErr := apperror.From(err)
	metrics.GetCollector().RecordError(c.Request.Context(), appErr)

	logFailure(c, appErr)

	c.SSEvent("error", gin.H{
		"status_code": appErr.Code,
		"status_msg":  appErr.Code.Msg(),
		"request_id":  RequestID(c),
		"message":     appErr.Code.Msg(),
		"error":       buildErrorDetail(appErr),
	})
}

func RequestID(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}
	return requestmeta.String(c.Request.Context(), requestmeta.FieldRequestID)
}

func logFailure(c *gin.Context, err error) {
	appErr := apperror.From(err)
	if appErr == nil {
		return
	}

	attrs := []any{
		"status_code", appErr.Code.Code(),
		"path", c.FullPath(),
		"method", c.Request.Method,
	}
	if appErr.Code.Code() >= code.CodeServerBusy.Code() {
		observe.Error(c.Request.Context(), "request failed", appErr, attrs...)
		return
	}

	observe.Warn(c.Request.Context(), "request handled with business error", append(attrs, observe.ErrorFields(appErr)...)...)
}

func buildErrorDetail(err error) *ErrorDetail {
	appErr := apperror.From(err)
	if appErr == nil {
		return nil
	}

	fields := apperror.FieldsOf(appErr)
	stringFields := make(map[string]string, len(fields))
	for key, value := range fields {
		if value == nil {
			continue
		}
		stringFields[key] = RequestValue(value)
	}
	if len(stringFields) == 0 {
		stringFields = nil
	}

	return &ErrorDetail{
		Code:    appErr.Code.Code(),
		Message: appErr.Message,
		Stack:   apperror.StackOf(appErr),
		Fields:  stringFields,
	}
}

func RequestValue(value any) string {
	return fmt.Sprintf("%v", value)
}
