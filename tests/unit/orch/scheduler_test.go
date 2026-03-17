package orch_test

import (
	"testing"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/orch"
)

func TestSchedulerBuildsGovernedUnits(t *testing.T) {
	scheduler := orch.NewScheduler()
	plan, err := scheduler.Build("sess_001", runcontract.RunContract{
		RunID:     "run_001",
		RiskClass: runcontract.RiskClassHigh,
		Metadata: map[string]string{
			"requires_shell":           "true",
			"requires_patch":           "true",
			"requires_tool_validation": "true",
		},
	}, "local_cli")
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if plan.Concurrency.MaxParallel != 1 {
		t.Fatalf("expected high risk to serialize work, got %d", plan.Concurrency.MaxParallel)
	}
	if len(plan.WorkUnits) < 5 {
		t.Fatalf("expected governed work units, got %d", len(plan.WorkUnits))
	}
}
