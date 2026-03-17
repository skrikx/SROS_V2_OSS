package runtime_integration_test

import (
	"context"
	"testing"
	"time"

	"srosv2/contracts/policy"
	"srosv2/internal/core/gov"
	"srosv2/internal/core/orch"
	"srosv2/internal/core/runtime"
)

func TestOrchGovRunSmoke(t *testing.T) {
	manager := newManager(t, &policy.Bundle{
		BundleID:              "pb_allow",
		Name:                  "Allow",
		Version:               "1",
		RulesetDigest:         "sha256:allow",
		DefaultVerdict:        policy.VerdictAllow,
		DefaultSandboxProfile: "local-default",
		Sandboxes: map[string]policy.SandboxProfile{
			"local-default": {Name: "local-default"},
		},
		Capabilities: []policy.CapabilityPolicy{
			{Name: "tool.validate", Verdict: policy.VerdictAllow, SandboxProfile: "local-default"},
		},
	})
	resp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: writeRuntimeContract(t, t.TempDir(), map[string]string{
		"requires_tool_validation": "true",
	})})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if resp.Session.State != runtime.SessionStateRunning {
		t.Fatalf("expected running state, got %s", resp.Session.State)
	}
	if resp.Decision != "allow" {
		t.Fatalf("expected allow decision, got %s", resp.Decision)
	}
}

func newManager(t *testing.T, bundle *policy.Bundle) *runtime.Manager {
	t.Helper()
	now := func() time.Time { return time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC) }
	root := t.TempDir()
	governor, err := gov.NewEngine(gov.Options{ArtifactRoot: root, Now: now, Bundle: bundle})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	orchestrator, err := orch.New(orch.Options{ArtifactRoot: root, Now: now})
	if err != nil {
		t.Fatalf("new orchestrator: %v", err)
	}
	manager, err := runtime.NewManager(runtime.Options{
		StoreDir:     root,
		Mode:         "local_cli",
		Now:          now,
		Gate:         stubGate{decision: runtime.AdmissionDecision{InitialState: runtime.SessionStateApproved, Reason: "allow", AutoStart: true}},
		Orchestrator: orchestrator,
		Governor:     governor,
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}
	return manager
}
