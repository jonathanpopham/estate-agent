package intake

import (
	"testing"

	"github.com/jonathanpopham/estate-agent/internal/work"
)

func TestErrorToWorkItem(t *testing.T) {
	payload := ErrorPayload{
		Service:     "checkout",
		Environment: "prod",
		Severity:    "high",
		Error:       "nil pointer",
		Stack:       "checkout.go:42",
	}

	first := ErrorToWorkItem(payload)
	second := ErrorToWorkItem(payload)

	if first.Kind != work.KindBug {
		t.Fatalf("Kind = %q, want %q", first.Kind, work.KindBug)
	}
	if first.ID != second.ID {
		t.Fatalf("ID not stable: %q != %q", first.ID, second.ID)
	}
	if first.Severity != "high" {
		t.Fatalf("Severity = %q", first.Severity)
	}
}
