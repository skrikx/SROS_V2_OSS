package trace_test

import (
	"testing"
	"time"

	ctrace "srosv2/contracts/trace"
	coretrace "srosv2/internal/core/trace"
)

func TestLedgerAppendOnly(t *testing.T) {
	service, err := coretrace.New(t.TempDir(), func() time.Time { return fixedTraceNow })
	if err != nil {
		t.Fatalf("new trace service: %v", err)
	}
	if _, err := service.Emit("run_001", "trace_001", "", "", ctrace.EventTypeRunStarted, map[string]any{"state": "started"}); err != nil {
		t.Fatalf("emit first event: %v", err)
	}
	if _, err := service.Emit("run_001", "trace_001", "", "", ctrace.EventTypeStateTransition, map[string]any{"to": "running"}); err != nil {
		t.Fatalf("emit second event: %v", err)
	}
	events, err := service.Reader.Events("run_001")
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 appended events, got %d", len(events))
	}
}
