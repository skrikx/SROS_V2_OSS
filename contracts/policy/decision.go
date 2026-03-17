package policy

import (
	"time"

	"srosv2/internal/shared/ids"
)

type Verdict string

const (
	VerdictAllow Verdict = "allow"
	VerdictAsk   Verdict = "ask"
	VerdictDeny  Verdict = "deny"
)

type PolicyDecision struct {
	ContractVersion      string               `json:"contract_version"`
	DecisionID           ids.PolicyDecisionID `json:"decision_id"`
	RunID                ids.RunID            `json:"run_id"`
	TraceID              ids.TraceID          `json:"trace_id"`
	Verdict              Verdict              `json:"verdict"`
	Boundary             TrustBoundary        `json:"boundary"`
	SandboxProfile       string               `json:"sandbox_profile"`
	ApprovalCheckpointID ids.CheckpointID     `json:"approval_checkpoint_id,omitempty"`
	BundleRef            ids.PolicyBundleID   `json:"bundle_ref"`
	Reason               string               `json:"reason"`
	EvidenceRefs         []ids.ArtifactID     `json:"evidence_refs,omitempty"`
	DecidedAt            time.Time            `json:"decided_at"`
}
