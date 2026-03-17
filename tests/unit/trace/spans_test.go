package trace_test

import (
	"testing"

	coretrace "srosv2/internal/core/trace"
)

func TestRootSpan(t *testing.T) {
	span := coretrace.RootSpan("run_001", "trace_001", "runtime.run")
	if span.SpanID == "" || span.TraceID == "" {
		t.Fatalf("unexpected root span: %+v", span)
	}
}
