package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"srosv2/contracts/runcontract"
)

type LifecycleEvent struct {
	From   SessionState `json:"from"`
	To     SessionState `json:"to"`
	Reason string       `json:"reason"`
	At     time.Time    `json:"at"`
}

type RuntimeSession struct {
	SessionID          string                  `json:"session_id"`
	RunID              string                  `json:"run_id"`
	ContractPath       string                  `json:"contract_path"`
	Contract           runcontract.RunContract `json:"contract"`
	State              SessionState            `json:"state"`
	Reason             string                  `json:"reason"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
	History            []LifecycleEvent        `json:"history"`
	ApprovalPath       string                  `json:"approval_path,omitempty"`
	LatestCheckpointID string                  `json:"latest_checkpoint_id,omitempty"`
	LatestRollbackID   string                  `json:"latest_rollback_id,omitempty"`
	PlanPath           string                  `json:"plan_path,omitempty"`
	LastDecision       string                  `json:"last_decision,omitempty"`
	TopologyBinding    string                  `json:"topology_binding,omitempty"`
	LatestMutationID   string                  `json:"latest_mutation_id,omitempty"`
	LatestWitnessID    string                  `json:"latest_witness_id,omitempty"`
}

func NewSession(contract runcontract.RunContract, contractPath string, now time.Time) RuntimeSession {
	sessionID := "sess_" + shortHash(string(contract.RunID)+"|"+now.UTC().Format(time.RFC3339Nano))
	return RuntimeSession{
		SessionID:    sessionID,
		RunID:        string(contract.RunID),
		ContractPath: contractPath,
		Contract:     contract,
		State:        SessionStatePlanned,
		Reason:       "runtime session created",
		CreatedAt:    now.UTC(),
		UpdatedAt:    now.UTC(),
		History:      []LifecycleEvent{},
	}
}

func RefFromSession(session RuntimeSession) SessionRef {
	return SessionRef{
		RunID:     session.RunID,
		SessionID: session.SessionID,
		State:     session.State,
		UpdatedAt: session.UpdatedAt,
	}
}

func shortHash(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])[:12]
}
