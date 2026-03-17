package trace_test

import (
	"testing"
	"time"

	ctrace "srosv2/contracts/trace"
	coretrace "srosv2/internal/core/trace"
)

func TestReaderEvents(t *testing.T) {
	service, err := coretrace.New(t.TempDir(), func() time.Time { return fixedTraceNow })
	if err != nil {
		t.Fatalf("new trace service: %v", err)
	}
	if _, err := service.Emit("run_001", "trace_001", "", "", ctrace.EventTypeRunStarted, nil); err != nil {
		t.Fatalf("emit event: %v", err)
	}
	events, err := service.Reader.Events("run_001")
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}
