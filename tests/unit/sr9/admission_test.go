package sr9_test

import (
	"testing"

	"srosv2/contracts/policy"
	"srosv2/contracts/runcontract"
	"srosv2/internal/core/sr9"
)

type stubPlanner struct {
	verdict policy.Verdict
	reason  string
}

func (s stubPlanner) Decide(_ runcontract.RunContract, _ sr9.Binding) (policy.Verdict, string) {
	return s.verdict, s.reason
}

func TestBuildAdmissionUsesPlannerVerdict(t *testing.T) {
	contract := validContract()
	admission, err := sr9.BuildAdmission(contract, stubPlanner{
		verdict: policy.VerdictDeny,
		reason:  "explicit deny for test",
	})
	if err != nil {
		t.Fatalf("build admission: %v", err)
	}
	if admission.Verdict != policy.VerdictDeny {
		t.Fatalf("expected deny verdict, got %s", admission.Verdict)
	}
}

func TestBuildAdmissionRejectsInvalidContract(t *testing.T) {
	contract := validContract()
	delete(contract.Metadata, "compile_request_id")
	if _, err := sr9.BuildAdmission(contract, nil); err == nil {
		t.Fatal("expected validation error")
	}
}
