package trace_test

import (
	"testing"
	"time"

	ctrace "srosv2/contracts/trace"
	coretrace "srosv2/internal/core/trace"
)

func TestReplayFromLineage(t *testing.T) {
	service, err := coretrace.New(t.TempDir(), func() time.Time { return fixedTraceNow })
	if err != nil {
		t.Fatalf("new trace service: %v", err)
	}
	_, _ = service.Emit("run_001", "trace_001", "", "", ctrace.EventTypeStateTransition, map[string]any{"to": "approved"})
	_, _ = service.Emit("run_001", "trace_001", "", "", ctrace.EventTypeStateTransition, map[string]any{"to": "running"})
	result, err := service.Replay.Replay("run_001")
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if result.EventCount != 2 || len(result.States) != 2 {
		t.Fatalf("unexpected replay result: %+v", result)
	}
}
