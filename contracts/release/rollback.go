package release

import (
	"time"

	"srosv2/internal/shared/ids"
)

type RollbackRecord struct {
	ContractVersion    string           `json:"contract_version"`
	RollbackID         ids.RollbackID   `json:"rollback_id"`
	ReleaseID          ids.ReleaseID    `json:"release_id"`
	TargetCheckpointID ids.CheckpointID `json:"target_checkpoint_id"`
	Reason             string           `json:"reason"`
	CreatedAt          time.Time        `json:"created_at"`
}
