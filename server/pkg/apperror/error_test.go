package apperror

import (
	"errors"
	"testing"

	"server/pkg/code"
)

func TestWrapCapturesCodeStackAndFields(t *testing.T) {
	root := errors.New("db timeout")
	err := Wrap(code.CodeServerBusy, root, "query failed").WithField("session_id", "session-1")

	if err == nil {
		t.Fatal("expected error")
	}

	if got := CodeOf(err); got != code.CodeServerBusy {
		t.Fatalf("unexpected code: got %v want %v", got, code.CodeServerBusy)
	}

	if got := MessageOf(err); got != "query failed" {
		t.Fatalf("unexpected message: got %q", got)
	}

	fields := FieldsOf(err)
	if got := fields["session_id"]; got != "session-1" {
		t.Fatalf("unexpected field value: got %v", got)
	}

	if len(StackOf(err)) == 0 {
		t.Fatal("expected stack frames to be captured")
	}
}

func TestFromWrapsGenericErrorAsServerBusy(t *testing.T) {
	err := From(errors.New("plain"))

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Code != code.CodeServerBusy {
		t.Fatalf("unexpected code: got %v want %v", err.Code, code.CodeServerBusy)
	}

	if len(err.Stack) == 0 {
		t.Fatal("expected stack frames")
	}
}
