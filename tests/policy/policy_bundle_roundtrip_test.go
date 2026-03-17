package policy_test

import (
	"testing"

	"srosv2/internal/core/gov"
)

func TestPolicyBundleRoundtrip(t *testing.T) {
	bundle, err := gov.LoadBundle("../../contracts/policy/local.bundle.json")
	if err != nil {
		t.Fatalf("load bundle: %v", err)
	}
	if string(bundle.BundleID) == "" {
		t.Fatalf("expected bundle id")
	}
}
