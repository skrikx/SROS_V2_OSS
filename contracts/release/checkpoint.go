package release

import (
	"time"

	"srosv2/internal/shared/ids"
)

type CheckpointRecord struct {
	ContractVersion string           `json:"contract_version"`
	CheckpointID    ids.CheckpointID `json:"checkpoint_id"`
	RunID           ids.RunID        `json:"run_id"`
	Stage           Stage            `json:"stage"`
	ArtifactRefs    []ids.ArtifactID `json:"artifact_refs,omitempty"`
	RecordedAt      time.Time        `json:"recorded_at"`
}
