package orch_test

import (
	"context"
	"testing"
	"time"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/orch"
)

func TestOrchestratorHydratesAndRoutesAsk(t *testing.T) {
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)
	o, err := orch.New(orch.Options{
		ArtifactRoot: t.TempDir(),
		Now:          func() time.Time { return now },
	})
	if err != nil {
		t.Fatalf("new orchestrator: %v", err)
	}

	plan, err := o.Hydrate("sess_001", runcontract.RunContract{
		RunID:          "run_001",
		RiskClass:      runcontract.RiskClassMedium,
		CheckpointRefs: []runcontract.CheckpointReference{{CheckpointID: "cp_001", Stage: "validated"}},
		Metadata: map[string]string{
			"requires_connector": "true",
		},
	}, "local_cli")
	if err != nil {
		t.Fatalf("hydrate: %v", err)
	}

	result, err := o.Execute(context.Background(), plan, func(_ context.Context, unit orch.WorkUnit) (orch.Decision, error) {
		if unit.Capability == "connector.invoke" {
			return orch.Decision{Verdict: "ask", Reason: "operator checkpoint required", SandboxProfile: "net-observe"}, nil
		}
		return orch.Decision{Verdict: "allow", Reason: "allowed"}, nil
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Completed {
		t.Fatal("expected execution to pause for ask")
	}
	if result.Route == nil || result.Route.ApprovalPath == "" {
		t.Fatal("expected checkpoint route with approval path")
	}
}
