package fabric_test

import (
	"testing"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/registry"
)

func TestLifecycleTransitionRules(t *testing.T) {
	if !registry.CanTransition(ctools.StateDraft, ctools.StateValidated) {
		t.Fatalf("expected draft -> validated to be allowed")
	}
	if registry.CanTransition(ctools.StateDraft, ctools.StateActive) {
		t.Fatalf("expected draft -> active to be blocked")
	}
}
