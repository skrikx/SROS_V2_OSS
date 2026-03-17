package orch_test

import (
	"testing"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/orch"
)

func TestConcurrencyRulesForRisk(t *testing.T) {
	low := orch.RulesForRisk(runcontract.RiskClassLow)
	high := orch.RulesForRisk(runcontract.RiskClassHigh)
	if low.MaxParallel <= high.MaxParallel {
		t.Fatalf("expected low risk to allow more parallelism than high risk")
	}
	if high.AllowSpawn {
		t.Fatal("expected high risk to disallow spawning")
	}
}
