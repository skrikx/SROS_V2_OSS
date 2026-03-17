package runcontract

import "srosv2/internal/shared/ids"

type CheckpointReference struct {
	CheckpointID ids.CheckpointID `json:"checkpoint_id"`
	Stage        string           `json:"stage"`
	Reason       string           `json:"reason,omitempty"`
}
