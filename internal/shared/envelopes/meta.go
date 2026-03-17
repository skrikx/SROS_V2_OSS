package envelopes

import (
	"time"

	"srosv2/internal/shared/ids"
)

// Meta carries neutral correlation and version metadata.
type Meta struct {
	ContractVersion string         `json:"contract_version"`
	RunID           ids.RunID      `json:"run_id,omitempty"`
	TraceID         ids.TraceID    `json:"trace_id,omitempty"`
	GeneratedAt     time.Time      `json:"generated_at"`
	Labels          []string       `json:"labels,omitempty"`
	Attributes      map[string]any `json:"attributes,omitempty"`
}
