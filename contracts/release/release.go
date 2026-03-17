package release

import (
	"time"

	"srosv2/internal/shared/ids"
)

type ReleaseRecord struct {
	ContractVersion string            `json:"contract_version"`
	ReleaseID       ids.ReleaseID     `json:"release_id"`
	CheckpointID    ids.CheckpointID  `json:"checkpoint_id"`
	TargetStage     Stage             `json:"target_stage"`
	PromotionMeta   map[string]string `json:"promotion_meta,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}
