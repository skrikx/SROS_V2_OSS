package fabric_test

import (
	"path/filepath"
	"testing"

	"srosv2/internal/fabric/harness"
	"srosv2/internal/fabric/mcpclient"
	"srosv2/internal/fabric/registry"
)

func TestMCPIngressNormalizesAndRegisters(t *testing.T) {
	reg, err := registry.New(filepath.Join(t.TempDir(), "registry"), harness.New(), nil)
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}
	manifest, err := mcpclient.Ingest("../../../examples/mcp/ingested_remote_capability.json", reg)
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if manifest.Name != "remote.docs.lookup" || !manifest.MCPIngressCapable {
		t.Fatalf("unexpected manifest: %+v", manifest)
	}
}
