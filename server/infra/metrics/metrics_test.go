package metrics

import (
	"context"
	"testing"

	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/requestmeta"
)

func TestCollectorRecordErrorCreatesAlertForServerErrors(t *testing.T) {
	collector := NewCollector()
	ctx := requestmeta.WithFields(context.Background(), requestmeta.Fields{
		requestmeta.FieldRequestID: "req-1",
		requestmeta.FieldUserName:  "alice",
		requestmeta.FieldSessionID: "session-1",
	})

	for range 3 {
		collector.RecordError(ctx, apperror.New(code.CodeServerBusy, "service busy").WithField("model_type", "qwen"))
	}

	errors := collector.RecentErrors()
	if len(errors) != 3 {
		t.Fatalf("unexpected recent error count: got %d want 3", len(errors))
	}

	if errors[0].RequestID != "req-1" {
		t.Fatalf("unexpected request id: got %q", errors[0].RequestID)
	}

	if len(errors[0].Stack) == 0 {
		t.Fatal("expected stack in recent error")
	}

	alerts := collector.AlertSnapshots()
	if len(alerts) != 1 {
		t.Fatalf("unexpected alert count: got %d want 1", len(alerts))
	}

	if alerts[0].Code != code.CodeServerBusy.Code() {
		t.Fatalf("unexpected alert code: got %d want %d", alerts[0].Code, code.CodeServerBusy.Code())
	}
}

func TestCollectorDoesNotAlertForBusinessErrors(t *testing.T) {
	collector := NewCollector()

	for range 3 {
		collector.RecordError(context.Background(), apperror.New(code.CodeInvalidParams, "invalid parameters"))
	}

	if len(collector.AlertSnapshots()) != 0 {
		t.Fatal("expected no alerts for business errors")
	}
}
