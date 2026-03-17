package rollback_test

import (
	"testing"
	"time"

	"srosv2/contracts/release"
)

func TestRollbackRecordValidation(t *testing.T) {
	errs := release.ValidateRollback(release.RollbackRecord{
		ContractVersion:    "v2.0",
		RollbackID:         "rb_001",
		ReleaseID:          "rel_001",
		TargetCheckpointID: "cp_001",
		Reason:             "operator rollback",
		CreatedAt:          time.Now().UTC(),
	})
	if len(errs) != 0 {
		t.Fatalf("expected valid rollback record, got %v", errs)
	}
}
