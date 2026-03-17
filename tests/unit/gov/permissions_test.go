package gov_test

import (
	"context"
	"testing"
	"time"

	"srosv2/contracts/policy"
	"srosv2/internal/core/gov"
)

func TestPatchDeniedWithoutExplicitPolicy(t *testing.T) {
	engine, err := gov.NewEngine(gov.Options{
		ArtifactRoot: t.TempDir(),
		Now:          func() time.Time { return time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC) },
		Bundle: &policy.Bundle{
			BundleID:              "pb_001",
			Name:                  "Patch Deny",
			Version:               "1",
			RulesetDigest:         "sha256:deny",
			DefaultVerdict:        policy.VerdictAllow,
			DefaultSandboxProfile: "local-default",
			Sandboxes: map[string]policy.SandboxProfile{
				"local-default": {Name: "local-default"},
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
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if result.Decision.Verdict != policy.VerdictDeny {
		t.Fatalf("expected deny, got %s", result.Decision.Verdict)
	}
}
