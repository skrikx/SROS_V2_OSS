package runtime_integration_test

import (
	"context"
	"path/filepath"
	"testing"

	"srosv2/contracts/policy"
	"srosv2/internal/core/gov"
)

func TestHighRiskCapabilityDenied(t *testing.T) {
	engine, err := gov.NewEngine(gov.Options{
		ArtifactRoot: filepath.Join(t.TempDir(), "artifacts"),
		Bundle: &policy.Bundle{
			BundleID:              "pb_001",
			Name:                  "deny shell",
			Version:               "1",
			RulesetDigest:         "sha256:deny-shell",
			DefaultSandboxProfile: "read_only",
			Sandboxes: map[string]policy.SandboxProfile{
				"read_only": {Name: "read_only"},
			},
		},
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	result, err := engine.Evaluate(context.Background(), gov.Request{
		RunID:      "run_high_risk",
		TraceID:    "trace_high_risk",
		Capability: "shell.exec",
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if result.Decision.Verdict != policy.VerdictDeny {
		t.Fatalf("expected deny, got %s", result.Decision.Verdict)
	}
}
