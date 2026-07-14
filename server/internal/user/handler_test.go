package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeUserService struct {
	loginIdentifier string
}

func (f *fakeUserService) Login(ctx context.Context, identifier, rawPassword string) (string, bool, error) {
	f.loginIdentifier = identifier
	return "token", false, nil
}

func (f *fakeUserService) Register(ctx context.Context, email, rawPassword, captcha string) (string, bool, error) {
	return "", false, nil
}

func (f *fakeUserService) SendCaptcha(ctx context.Context, email string) error {
	return nil
}

func TestLoginUsesEmailField(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeUserService{}
	handler := NewHandler(service)
	router := gin.New()
	router.POST("/login", handler.Login)

	body, _ := json.Marshal(map[string]string{
		"email":    "User@qq.com",
		"password": "secret",
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if service.loginIdentifier != "User@qq.com" {
		t.Fatalf("unexpected login identifier: got %q", service.loginIdentifier)
	}
}

func TestLoginRejectsLegacyUsernameField(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeUserService{}
	handler := NewHandler(service)
	router := gin.New()
	router.POST("/login", handler.Login)

	body, _ := json.Marshal(map[string]string{
		"username": "legacy@example.com",
		"password": "secret",
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if service.loginIdentifier != "" {
		t.Fatalf("unexpected login identifier: got %q", service.loginIdentifier)
	}

	var payload map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if payload["status_code"] != float64(2001) {
		t.Fatalf("unexpected status code: got %v want %v", payload["status_code"], 2001)
	}
}
