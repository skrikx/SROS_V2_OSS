package memory

import (
	"time"

	"srosv2/internal/shared/ids"
)

type BranchReference struct {
	BranchID       ids.BranchID         `json:"branch_id"`
	ParentBranchID ids.BranchID         `json:"parent_branch_id,omitempty"`
	HeadMutationID ids.MemoryMutationID `json:"head_mutation_id,omitempty"`
}

type BranchRecord struct {
	BranchID       ids.BranchID         `json:"branch_id"`
	ParentBranchID ids.BranchID         `json:"parent_branch_id,omitempty"`
	HeadMutationID ids.MemoryMutationID `json:"head_mutation_id,omitempty"`
	CreatedBy      ids.OperatorID       `json:"created_by"`
	TenantID       ids.TenantID         `json:"tenant_id"`
	WorkspaceID    ids.WorkspaceID      `json:"workspace_id"`
	CreatedAt      time.Time            `json:"created_at"`
	RewoundTo      ids.MemoryMutationID `json:"rewound_to,omitempty"`
	Reason         string               `json:"reason,omitempty"`
}
