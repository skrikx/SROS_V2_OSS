package memory

import (
	"time"

	"srosv2/internal/shared/ids"
)

type MutationKind string

const (
	MutationKindUpsert           MutationKind = "upsert"
	MutationKindDelete           MutationKind = "delete"
	MutationKindLink             MutationKind = "link"
	MutationKindAnnotate         MutationKind = "annotate"
	MutationKindPruneRecommend   MutationKind = "prune_recommend"
	MutationKindCompactRecommend MutationKind = "compact_recommend"
)

type MemoryMutation struct {
	ContractVersion  string               `json:"contract_version"`
	MutationID       ids.MemoryMutationID `json:"mutation_id"`
	RunID            ids.RunID            `json:"run_id"`
	SessionID        ids.SessionID        `json:"session_id"`
	TenantID         ids.TenantID         `json:"tenant_id"`
	WorkspaceID      ids.WorkspaceID      `json:"workspace_id"`
	OperatorID       ids.OperatorID       `json:"operator_id"`
	Scope            Scope                `json:"scope"`
	Kind             MutationKind         `json:"kind"`
	Branch           BranchReference      `json:"branch"`
	ParentMutationID ids.MemoryMutationID `json:"parent_mutation_id,omitempty"`
	LineageRef       string               `json:"lineage_ref"`
	RecallIndexRef   string               `json:"recall_index_ref,omitempty"`
	Key              string               `json:"key"`
	Value            string               `json:"value,omitempty"`
	Reason           string               `json:"reason,omitempty"`
	OccurredAt       time.Time            `json:"occurred_at"`
}
