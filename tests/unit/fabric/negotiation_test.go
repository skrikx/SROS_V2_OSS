package fabric_test

import (
	"path/filepath"
	"testing"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/harness"
	"srosv2/internal/fabric/mcpclient"
	"srosv2/internal/fabric/registry"
)

func TestNegotiationSelectsRunnableCapability(t *testing.T) {
	reg, err := registry.New(filepath.Join(t.TempDir(), "registry"), harness.New(), []ctools.Manifest{
		{Name: "local.patch", Title: "Patch", Description: "Patch", ManifestVersion: "v2.0", Version: "1", Class: "tool.local.patch", Domain: "workspace", PolicyClass: "patch.apply", Status: ctools.StateActive, TrustBoundary: "local_fs", SandboxProfiles: []string{"patch_only"}},
	})
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}
	result := mcpclient.Negotiate(reg, ctools.NegotiationRequest{Query: ctools.SearchQuery{Class: "tool.local"}, RequireRunnable: true})
	if !result.Allowed || result.SelectedCapability != "local.patch" {
		t.Fatalf("unexpected negotiation result: %+v", result)
	}
}
