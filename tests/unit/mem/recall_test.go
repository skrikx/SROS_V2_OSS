package mem_test

import (
	"testing"
	"time"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/core/mem"
)

func TestRecallQuery(t *testing.T) {
	store, err := mem.NewStore(t.TempDir(), func() time.Time { return fixedMemNow })
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	if _, err := store.Upsert(mem.MutationInput{
		Scope:      mem.ScopeBinding{Scope: cmemory.ScopeWorkspace, TenantID: "local", WorkspaceID: "default", RunID: "run_001", SessionID: "sess_001"},
		OperatorID: "op_local",
		Kind:       cmemory.MutationKindUpsert,
		Branch:     cmemory.BranchReference{BranchID: "branch_main"},
		Key:        "workspace.intent",
		Value:      "workspace compile summary",
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	results, err := store.Recall("workspace")
	if err != nil {
		t.Fatalf("recall: %v", err)
	}
	if results["match_count"].(int) == 0 {
		t.Fatal("expected recall results")
	}
}
