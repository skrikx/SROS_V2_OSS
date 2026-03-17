package mem_test

import (
	"testing"
	"time"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/core/mem"
)

func TestBranchAndRewind(t *testing.T) {
	store, err := mem.NewStore(t.TempDir(), func() time.Time { return fixedMemNow })
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	first, err := store.Upsert(mem.MutationInput{
		Scope:      mem.ScopeBinding{Scope: cmemory.ScopeWorkspace, TenantID: "local", WorkspaceID: "default", RunID: "run_001", SessionID: "sess_001"},
		OperatorID: "op_local",
		Kind:       cmemory.MutationKindUpsert,
		Branch:     cmemory.BranchReference{BranchID: "branch_main"},
		Key:        "workspace.intent",
		Value:      "seed",
		OccurredAt: fixedMemNow,
	})
	if err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	if err := store.Rewind("branch_main", first.MutationID, "op_local", "test rewind"); err != nil {
		t.Fatalf("rewind: %v", err)
	}
	branches, err := store.Branches()
	if err != nil {
		t.Fatalf("branches: %v", err)
	}
	if len(branches) == 0 || branches[0].RewoundTo != first.MutationID {
		t.Fatalf("expected rewound branch, got %+v", branches)
	}
}
