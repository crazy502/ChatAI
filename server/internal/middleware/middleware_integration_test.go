package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/pkg/code"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

type okPayload struct {
	Value string `json:"value"`
	response.Response
}

type errorEnvelope struct {
	StatusCode int64  `json:"status_code"`
	RequestID  string `json:"request_id"`
	Error      struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestRequestContextPropagatesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestContext(), Recovery(), RequestObserver())
	router.GET("/ok", func(c *gin.Context) {
		response.OK(c, &okPayload{Value: "pong"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req.Header.Set("X-Request-ID", "req-fixed")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Header().Get("X-Request-ID") != "req-fixed" {
		t.Fatalf("unexpected response request id header: got %q", recorder.Header().Get("X-Request-ID"))
	}

	var payload okPayload
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if payload.RequestID != "req-fixed" {
		t.Fatalf("unexpected request id in body: got %q", payload.RequestID)
	}
}

func TestRecoveryReturnsStructuredError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestContext(), Recovery(), RequestObserver())
	router.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected http status: got %d want %d", recorder.Code, http.StatusInternalServerError)
	}

	var payload errorEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if payload.StatusCode != code.CodeServerBusy.Code() {
		t.Fatalf("unexpected status code: got %d want %d", payload.StatusCode, code.CodeServerBusy.Code())
	}

	if payload.RequestID == "" {
		t.Fatal("expected request id in error response")
	}

	if payload.Error.Code != code.CodeServerBusy.Code() {
		t.Fatalf("unexpected error code: got %d want %d", payload.Error.Code, code.CodeServerBusy.Code())
	}

	if payload.Error.Message != code.CodeServerBusy.Msg() {
		t.Fatalf("unexpected public error message: got %q want %q", payload.Error.Message, code.CodeServerBusy.Msg())
	}

	var rawPayload map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &rawPayload); err != nil {
		t.Fatalf("unmarshal raw response: %v", err)
	}
	errorDetail, ok := rawPayload["error"].(map[string]any)
	if !ok {
		t.Fatal("expected error detail")
	}
	if _, exists := errorDetail["stack"]; exists {
		t.Fatal("stack must not be exposed in public error payload")
	}
	if _, exists := errorDetail["fields"]; exists {
		t.Fatal("internal fields must not be exposed in public error payload")
	}
}
