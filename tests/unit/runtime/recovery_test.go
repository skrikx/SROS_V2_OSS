package runtime_test

import (
	"context"
	"testing"
	"time"

	"srosv2/internal/core/runtime"
)

func TestManagerRollbackUsesCheckpointRecord(t *testing.T) {
	dir := t.TempDir()
	contractPath := writeRunContractFile(t, dir, nil)

	manager, err := runtime.NewManager(runtime.Options{
		StoreDir: dir,
		Now:      func() time.Time { return fixedNow },
		Gate: stubGate{decision: runtime.AdmissionDecision{
			InitialState: runtime.SessionStateApproved,
			Reason:       "allow",
			AutoStart:    true,
		}},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	runResp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: contractPath})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	cpResp, err := manager.Checkpoint(context.Background(), runtime.CheckpointRequest{
		SessionID: runResp.Session.SessionID,
		Stage:     "validated",
	})
	if err != nil {
		t.Fatalf("checkpoint: %v", err)
	}

	rbResp, err := manager.Rollback(context.Background(), runtime.RollbackRequest{
		SessionID:    runResp.Session.SessionID,
		CheckpointID: cpResp.CheckpointID,
		Reason:       "operator rollback",
	})
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rbResp.RollbackID == "" {
		t.Fatal("expected rollback id")
	}
	if rbResp.Session.State != runtime.SessionStateRolledBack {
		t.Fatalf("expected rolled_back state, got %s", rbResp.Session.State)
	}
}
