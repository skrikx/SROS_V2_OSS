package trace_test

import (
	"testing"
	"time"

	ctrace "srosv2/contracts/trace"
	coretrace "srosv2/internal/core/trace"
)

func TestQueryByType(t *testing.T) {
	service, err := coretrace.New(t.TempDir(), func() time.Time { return fixedTraceNow })
	if err != nil {
		t.Fatalf("new trace service: %v", err)
	}
	_, _ = service.Emit("run_001", "trace_001", "", "", ctrace.EventTypeRunStarted, nil)
	_, _ = service.Emit("run_001", "trace_001", "", "", ctrace.EventTypeMemoryMutation, nil)
	events, err := service.Query.ByType("run_001", ctrace.EventTypeMemoryMutation)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 matching event, got %d", len(events))
	}
}
