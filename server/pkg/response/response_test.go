package response

import (
	"net/http"
	"testing"

	"server/pkg/code"
)

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		name       string
		resultCode code.Code
		want       int
	}{
		{name: "invalid params", resultCode: code.CodeInvalidParams, want: http.StatusBadRequest},
		{name: "invalid token", resultCode: code.CodeInvalidToken, want: http.StatusUnauthorized},
		{name: "forbidden", resultCode: code.CodeForbidden, want: http.StatusForbidden},
		{name: "not found", resultCode: code.CodeRecordNotFound, want: http.StatusNotFound},
		{name: "conflict", resultCode: code.CodeEmailExist, want: http.StatusConflict},
		{name: "rate limited", resultCode: code.CodeTooManyRequests, want: http.StatusTooManyRequests},
		{name: "server error", resultCode: code.CodeServerBusy, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httpStatus(tt.resultCode); got != tt.want {
				t.Fatalf("unexpected http status: got %d want %d", got, tt.want)
			}
		})
	}
}
