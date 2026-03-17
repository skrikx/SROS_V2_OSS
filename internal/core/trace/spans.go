package trace

import (
	ctrace "srosv2/contracts/trace"
	"srosv2/internal/shared/ids"
)

func RootSpan(runID ids.RunID, traceID ids.TraceID, name string) ctrace.Span {
	return ctrace.Span{
		TraceID:    traceID,
		SpanID:     ids.SpanID("span_root_" + shortHash(string(runID)+"|"+name)),
		Name:       name,
		Attributes: map[string]string{"run_id": string(runID)},
	}
}
