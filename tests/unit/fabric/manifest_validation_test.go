package fabric_test

import (
	"testing"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/harness"
	"srosv2/internal/fabric/registry"
)

func TestManifestValidationRequiresSandboxProfiles(t *testing.T) {
	result := registry.ValidateManifest(ctools.Manifest{
		ManifestVersion: "v2.0",
		Name:            "invalid",
		Title:           "Invalid",
		Description:     "missing sandbox",
		Version:         "1.0.0",
		Class:           "tool.local.patch",
		Domain:          "workspace",
		PolicyClass:     "patch.apply",
		Status:          ctools.StateDraft,
		TrustBoundary:   "local_fs",
	}, harness.New())
	if result.Valid {
		t.Fatalf("expected invalid result")
	}
}
