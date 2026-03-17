package intake_test

import (
	"testing"

	"srosv2/internal/ace/intake"
)

func TestBuildPlan(t *testing.T) {
	plan := intake.BuildPlan("research and analyze the incident")
	if plan.Classification.Domain != intake.DomainResearch {
		t.Fatalf("expected research domain, got %s", plan.Classification.Domain)
	}
	if len(plan.Shortlist.Skills) == 0 {
		t.Fatal("expected skill shortlist")
	}
	if len(plan.Preflight) == 0 {
		t.Fatal("expected preflight checks")
	}
}
