package runcontract

import (
	"time"

	"srosv2/internal/shared/ids"
)

type RiskClass string

const (
	RiskClassLow      RiskClass = "low"
	RiskClassMedium   RiskClass = "medium"
	RiskClassHigh     RiskClass = "high"
	RiskClassCritical RiskClass = "critical"
)

type RouteClass string

const (
	RouteClassLocalCLI RouteClass = "local_cli"
	RouteClassLocalRun RouteClass = "local_run"
)

type ReceiptMode string

const (
	ReceiptModeNone    ReceiptMode = "none"
	ReceiptModeSummary ReceiptMode = "summary"
	ReceiptModeFull    ReceiptMode = "full"
)

type ArtifactReference struct {
	ArtifactID ids.ArtifactID `json:"artifact_id"`
	Kind       string         `json:"kind"`
	URI        string         `json:"uri"`
}

type ReceiptRequest struct {
	Mode ReceiptMode `json:"mode"`
	Kind string      `json:"kind"`
}

type RunContract struct {
	ContractVersion   string                `json:"contract_version"`
	RunID             ids.RunID             `json:"run_id"`
	TraceID           ids.TraceID           `json:"trace_id"`
	OperatorID        ids.OperatorID        `json:"operator_id"`
	TenantID          ids.TenantID          `json:"tenant_id"`
	WorkspaceID       ids.WorkspaceID       `json:"workspace_id"`
	IntentSummary     string                `json:"intent_summary"`
	NormalizedRequest string                `json:"normalized_request"`
	RiskClass         RiskClass             `json:"risk_class"`
	RouteClass        RouteClass            `json:"route_class"`
	ArtifactRefs      []ArtifactReference   `json:"artifact_refs,omitempty"`
	CheckpointRefs    []CheckpointReference `json:"checkpoint_refs,omitempty"`
	RequestedReceipts []ReceiptRequest      `json:"requested_receipts,omitempty"`
	Metadata          map[string]string     `json:"metadata,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
}
