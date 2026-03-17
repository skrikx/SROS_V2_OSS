package runtime

import (
	"fmt"
	"strings"
	"time"

	"srosv2/contracts/release"
	"srosv2/internal/shared/ids"
)

type RuntimeRollback struct {
	SessionID string                 `json:"session_id"`
	Record    release.RollbackRecord `json:"record"`
}

func NewRollback(session RuntimeSession, checkpointID, reason string, now time.Time) (RuntimeRollback, error) {
	if strings.TrimSpace(checkpointID) == "" {
		return RuntimeRollback{}, fmt.Errorf("checkpoint id is required for rollback")
	}
	if strings.TrimSpace(reason) == "" {
		reason = "operator rollback"
	}

	record := release.RollbackRecord{
		ContractVersion:    "v2.0",
		RollbackID:         ids.RollbackID("rb_" + shortHash(session.SessionID+"|"+now.UTC().Format(time.RFC3339Nano))),
		ReleaseID:          ids.ReleaseID("rel_" + shortHash(checkpointID+"|release")),
		TargetCheckpointID: ids.CheckpointID(checkpointID),
		Reason:             reason,
		CreatedAt:          now.UTC(),
	}

	if errs := release.ValidateRollback(record); len(errs) > 0 {
		return RuntimeRollback{}, fmt.Errorf("invalid rollback record: %v", errs[0])
	}

	return RuntimeRollback{SessionID: session.SessionID, Record: record}, nil
}
