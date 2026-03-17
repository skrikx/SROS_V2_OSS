package trace

import (
	"time"

	"srosv2/internal/shared/ids"
)

type TraceEvent struct {
	ContractVersion string         `json:"contract_version"`
	EventID         ids.EventID    `json:"event_id"`
	TraceID         ids.TraceID    `json:"trace_id"`
	SpanID          ids.SpanID     `json:"span_id"`
	ParentSpanID    ids.SpanID     `json:"parent_span_id,omitempty"`
	RunID           ids.RunID      `json:"run_id"`
	EventType       EventType      `json:"event_type"`
	OccurredAt      time.Time      `json:"occurred_at"`
	Payload         map[string]any `json:"payload,omitempty"`
	ArtifactRefs    []ids.ArtifactID `json:"artifact_refs,omitempty"`
	ReceiptRef      ids.ReceiptID  `json:"receipt_ref,omitempty"`
}
