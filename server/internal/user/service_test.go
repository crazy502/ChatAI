package user

import "testing"

func TestNormalizeEmail(t *testing.T) {
	got := normalizeEmail("  User@Example.COM ")
	if got != "user@example.com" {
		t.Fatalf("unexpected normalized email: got %q", got)
	}
}
