package user

import "testing"

func TestNormalizeEmail(t *testing.T) {
	got := normalizeEmail("  User@Example.COM ")
	if got != "user@example.com" {
		t.Fatalf("unexpected normalized email: got %q", got)
	}
}

func TestNormalizeEmailAppendsQQSuffixWhenMissing(t *testing.T) {
	got := normalizeEmail(" 12345678 ")
	if got != "12345678@qq.com" {
		t.Fatalf("unexpected normalized email with suffix: got %q", got)
	}
}
