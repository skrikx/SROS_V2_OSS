package mem_test

import (
	"testing"
	"time"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/core/mem"
)

func TestStoreUpsertCreatesLineage(t *testing.T) {
	store, err := mem.NewStore(t.TempDir(), func() time.Time { return fixedMemNow })
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	mutation, err := store.Upsert(mem.MutationInput{
		Scope: mem.ScopeBinding{
			Scope:       cmemory.ScopeWorkspace,
			TenantID:    "local",
			WorkspaceID: "default",
			RunID:       "run_001",
			SessionID:   "sess_001",
		},
		OperatorID: "op_local",
		Kind:       cmemory.MutationKindUpsert,
		Branch:     cmemory.BranchReference{BranchID: "branch_main"},
		Key:        "workspace.intent",
		Value:      "remember compile state",
		Reason:     "test",
	})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if mutation.LineageRef == "" {
		t.Fatal("expected lineage ref")
	}
	records, err := store.Records()
	if err != nil {
		t.Fatalf("records: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
}
