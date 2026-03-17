package fabric_test

import (
	"testing"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/registry"
)

func TestRegistrySearchRanksActiveCapabilities(t *testing.T) {
	results := registry.Search([]ctools.Manifest{
		{Name: "local.shell", Title: "Shell", Class: "tool.local.shell", Domain: "workspace", PolicyClass: "shell.exec", Status: ctools.StateExperimental, TrustBoundary: "local_process", SandboxProfiles: []string{"shell_gated"}},
		{Name: "local.patch", Title: "Patch", Class: "tool.local.patch", Domain: "workspace", PolicyClass: "patch.apply", Status: ctools.StateActive, TrustBoundary: "local_fs", SandboxProfiles: []string{"patch_only"}},
	}, ctools.SearchQuery{Class: "tool.local"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Manifest.Name != "local.patch" {
		t.Fatalf("expected local.patch first, got %s", results[0].Manifest.Name)
	}
}
