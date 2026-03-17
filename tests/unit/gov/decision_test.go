package gov_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"srosv2/contracts/policy"
	"srosv2/internal/core/gov"
)

func TestEngineAllowDecisionMatchesGolden(t *testing.T) {
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)
	engine, err := gov.NewEngine(gov.Options{
		ArtifactRoot: t.TempDir(),
		Now:          func() time.Time { return now },
		Bundle: &policy.Bundle{
			BundleID:              "pb_001",
			Name:                  "Golden Allow",
			Version:               "1",
			RulesetDigest:         "sha256:allow",
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
		Capability: "runtime.session.prepare",
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if result.Decision.Verdict != policy.VerdictAllow {
		t.Fatalf("expected allow, got %s", result.Decision.Verdict)
	}

	path := filepath.Join("..", "..", "golden", "gov", "allow_decision.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	var golden policy.PolicyDecision
	if err := json.Unmarshal(data, &golden); err != nil {
		t.Fatalf("decode golden: %v", err)
	}
	if golden.Verdict != result.Decision.Verdict || golden.Boundary != result.Decision.Boundary {
		t.Fatalf("golden mismatch")
	}
}
