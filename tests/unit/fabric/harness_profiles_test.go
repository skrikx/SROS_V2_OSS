package fabric_test

import (
	"testing"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/harness"
)

func TestHarnessResolvesDeclaredProfile(t *testing.T) {
	resolution, err := harness.New().Resolve(ctools.Manifest{
		Name:            "local.patch",
		Class:           "tool.local.patch",
		SandboxProfiles: []string{"patch_only"},
	})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !resolution.Compatible || resolution.Profile.Name != "patch_only" {
		t.Fatalf("unexpected resolution: %+v", resolution)
	}
}
