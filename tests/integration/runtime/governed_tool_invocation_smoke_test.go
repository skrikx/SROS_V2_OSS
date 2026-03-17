package runtime_integration_test

import (
	"context"
	"path/filepath"
	"testing"

	"srosv2/internal/core/boot"
	"srosv2/internal/shared/config"
)

func TestGovernedToolInvocationSmoke(t *testing.T) {
	cfg := config.Config{
		Mode:             config.ModeLocalCLI,
		WorkspaceRoot:    t.TempDir(),
		ArtifactRoot:     filepath.Join(t.TempDir(), "artifacts"),
		PolicyBundlePath: filepath.Join("..", "..", "..", "examples", "policy", "local_default_policy.json"),
		OutputFormat:     "text",
	}
	bundle, err := boot.Bootstrap(cfg)
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	payload, err := bundle.Fabric.ConnectorsInspectEnvelope(context.Background(), "../../../examples/connectors/local_secret_envelope.json")
	if err != nil {
		t.Fatalf("inspect envelope: %v", err)
	}
	if payload["secret_material"] != "[REDACTED]" {
		t.Fatalf("expected redacted envelope output")
	}
}
