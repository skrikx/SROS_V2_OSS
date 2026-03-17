package mem_test

import (
	"testing"
	"time"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/core/mem"
)

func TestMutationLedgerAppend(t *testing.T) {
	store, err := mem.NewStore(t.TempDir(), func() time.Time { return fixedMemNow })
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	if _, err := store.Upsert(mem.MutationInput{
		Scope:      mem.ScopeBinding{Scope: cmemory.ScopeSession, TenantID: "local", WorkspaceID: "default", RunID: "run_001", SessionID: "sess_001"},
		OperatorID: "op_local",
		Kind:       cmemory.MutationKindAnnotate,
		Branch:     cmemory.BranchReference{BranchID: "branch_main"},
		Key:        "runtime.state",
		Value:      "running",
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	entries, err := store.Ledger()
	if err != nil {
		t.Fatalf("ledger: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}
