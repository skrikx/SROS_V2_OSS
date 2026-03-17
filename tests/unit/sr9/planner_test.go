package sr9_test

import (
	"testing"

	"srosv2/contracts/policy"
	"srosv2/contracts/runcontract"
	"srosv2/internal/core/sr9"
)

func TestDefaultPlannerAsksForCriticalRisk(t *testing.T) {
	contract := validContract()
	contract.RiskClass = runcontract.RiskClassCritical

	admission, err := sr9.BuildAdmission(contract, nil)
	if err != nil {
		t.Fatalf("build admission: %v", err)
	}
	if admission.Verdict != policy.VerdictAsk {
		t.Fatalf("expected ask verdict for critical risk, got %s", admission.Verdict)
	}
}

func TestDefaultPlannerRespectsDenyOverride(t *testing.T) {
	contract := validContract()
	contract.Metadata["approval_mode"] = "deny"

	admission, err := sr9.BuildAdmission(contract, nil)
	if err != nil {
		t.Fatalf("build admission: %v", err)
	}
	if admission.Verdict != policy.VerdictDeny {
		t.Fatalf("expected deny verdict, got %s", admission.Verdict)
	}
}
