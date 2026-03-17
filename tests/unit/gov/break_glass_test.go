package gov_test

import (
	"context"
	"testing"
	"time"

	"srosv2/contracts/policy"
	"srosv2/internal/core/gov"
)

func TestBreakGlassRequiresAsk(t *testing.T) {
	engine, err := gov.NewEngine(gov.Options{
		ArtifactRoot: t.TempDir(),
		Now:          func() time.Time { return time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC) },
		Bundle: &policy.Bundle{
			BundleID:              "pb_001",
			Name:                  "Break Glass",
			Version:               "1",
			RulesetDigest:         "sha256:bg",
			DefaultSandboxProfile: "patch-safe",
			BreakGlassAllowed:     true,
			Sandboxes: map[string]policy.SandboxProfile{
				"patch-safe": {Name: "patch-safe", AllowPatch: true},
			},
			Capabilities: []policy.CapabilityPolicy{
				{Name: "patch.apply", Verdict: policy.VerdictAllow, SandboxProfile: "patch-safe", AllowBreakGlass: true},
			},
		},
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	result, err := engine.Evaluate(context.Background(), gov.Request{
		RunID:      "run_001",
		TraceID:    "trace_001",
		Capability: "patch.apply",
		BreakGlass: true,
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if result.Decision.Verdict != policy.VerdictAsk {
		t.Fatalf("expected ask, got %s", result.Decision.Verdict)
	}
}
