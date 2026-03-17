package trace_test

import (
	"testing"
	"time"

	coretrace "srosv2/internal/core/trace"
)

func TestWriterBeginSpan(t *testing.T) {
	service, err := coretrace.New(t.TempDir(), func() time.Time { return fixedTraceNow })
	if err != nil {
		t.Fatalf("new trace service: %v", err)
	}
	span, err := service.Writer.BeginSpan("trace_001", "", "runtime.run", map[string]string{"run_id": "run_001"})
	if err != nil {
		t.Fatalf("begin span: %v", err)
	}
	if span.SpanID == "" {
		t.Fatal("expected span id")
	}
}
