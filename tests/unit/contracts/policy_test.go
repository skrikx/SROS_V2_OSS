package contracts_test

import (
	"encoding/json"
	"testing"
	"time"

	"srosv2/contracts/policy"
	"srosv2/internal/shared/ids"
)

func TestPolicyDecisionValidateValid(t *testing.T) {
	decision := policy.PolicyDecision{
		ContractVersion: "v2.0",
		DecisionID:      ids.PolicyDecisionID("pd_001"),
		RunID:           ids.RunID("run_001"),
		TraceID:         ids.TraceID("trace_001"),
		Verdict:         policy.VerdictAllow,
		Boundary:        policy.TrustBoundaryLocalProcess,
		SandboxProfile:  "local-default",
		BundleRef:       ids.PolicyBundleID("pb_001"),
		Reason:          "safe local operation",
		DecidedAt:       time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
	}

	errs := policy.ValidateDecision(decision)
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(errs))
	}
}

func TestPolicyDecisionGoldenFixture(t *testing.T) {
	data := loadFixture(t, "policy_decision.json")
	var decision policy.PolicyDecision
	if err := json.Unmarshal(data, &decision); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	errs := policy.ValidateDecision(decision)
	if len(errs) != 0 {
		t.Fatalf("expected fixture to validate, got %d errors", len(errs))
	}
}
