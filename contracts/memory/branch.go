package memory

import "srosv2/internal/shared/ids"

type BranchReference struct {
	BranchID       ids.BranchID `json:"branch_id"`
	ParentBranchID ids.BranchID `json:"parent_branch_id,omitempty"`
}
