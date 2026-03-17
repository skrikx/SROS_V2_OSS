package fabric_test

import (
	"testing"

	"srosv2/internal/fabric/harness"
)

func TestSandboxCompatibility(t *testing.T) {
	if !harness.Compatible(harness.Defaults()["patch_only"], "tool.local.patch") {
		t.Fatalf("expected patch profile to allow tool.local.patch")
	}
	if harness.Compatible(harness.Defaults()["read_only"], "tool.local.shell") {
		t.Fatalf("expected read_only to deny shell")
	}
}
