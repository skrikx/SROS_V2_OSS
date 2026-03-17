package contracts_test

import (
	"encoding/json"
	"testing"
	"time"

	"srosv2/contracts/memory"
	"srosv2/internal/shared/ids"
)

func TestMemoryMutationValidateValid(t *testing.T) {
	mutation := memory.MemoryMutation{
		ContractVersion: "v2.0",
		MutationID:      ids.MemoryMutationID("mm_001"),
		RunID:           ids.RunID("run_001"),
		SessionID:       ids.SessionID("session_001"),
		Scope:           memory.ScopeSession,
		Kind:            memory.MutationKindUpsert,
		Branch: memory.BranchReference{
			BranchID: ids.BranchID("branch_main"),
		},
		Key:        "intent.summary",
		OccurredAt: time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
	}

	errs := memory.ValidateMutation(mutation)
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(errs))
	}
}

func TestMemoryMutationGoldenFixture(t *testing.T) {
	data := loadFixture(t, "memory_mutation.json")
	var mutation memory.MemoryMutation
	if err := json.Unmarshal(data, &mutation); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	errs := memory.ValidateMutation(mutation)
	if len(errs) != 0 {
		t.Fatalf("expected fixture to validate, got %d errors", len(errs))
	}
}
