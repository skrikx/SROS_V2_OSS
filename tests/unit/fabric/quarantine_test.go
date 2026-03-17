package fabric_test

import (
	"testing"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/registry"
)

func TestQuarantinedCapabilitiesRemainVisibleButNotSelectable(t *testing.T) {
	results := registry.Search([]ctools.Manifest{
		{Name: "local.shell", Title: "Shell", Class: "tool.local.shell", Domain: "workspace", PolicyClass: "shell.exec", Status: ctools.StateQuarantined, TrustBoundary: "local_process", SandboxProfiles: []string{"shell_gated"}},
	}, ctools.SearchQuery{Class: "tool.local", IncludeHistorical: true})
	if len(results) != 1 {
		t.Fatalf("expected one result")
	}
	if results[0].Selectable {
		t.Fatalf("expected quarantined capability to be non-selectable")
	}
}
