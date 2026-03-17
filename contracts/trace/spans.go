package trace

import (
	"time"

	"srosv2/internal/shared/ids"
)

type Span struct {
	TraceID      ids.TraceID `json:"trace_id"`
	SpanID       ids.SpanID  `json:"span_id"`
	ParentSpanID ids.SpanID  `json:"parent_span_id,omitempty"`
	Name         string      `json:"name"`
	StartedAt    time.Time   `json:"started_at"`
	EndedAt      time.Time   `json:"ended_at,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}
