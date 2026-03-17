package contracts_test

import (
	"encoding/json"
	"testing"
	"time"

	ctrace "srosv2/contracts/trace"
	"srosv2/internal/shared/ids"
)

func TestTraceEventValidateValid(t *testing.T) {
	event := ctrace.TraceEvent{
		ContractVersion: "v2.0",
		EventID:         ids.EventID("event_001"),
		TraceID:         ids.TraceID("trace_001"),
		SpanID:          ids.SpanID("span_001"),
		RunID:           ids.RunID("run_001"),
		EventType:       ctrace.EventTypeRunStarted,
		OccurredAt:      time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
	}

	errs := ctrace.ValidateEvent(event)
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(errs))
	}
}

func TestTraceEventGoldenFixture(t *testing.T) {
	data := loadFixture(t, "trace_event.json")
	var event ctrace.TraceEvent
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	errs := ctrace.ValidateEvent(event)
	if len(errs) != 0 {
		t.Fatalf("expected fixture to validate, got %d errors", len(errs))
	}
}
