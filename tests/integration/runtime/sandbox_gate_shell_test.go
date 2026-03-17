package runtime_integration_test

import (
	"context"
	"testing"

	"srosv2/contracts/policy"
	"srosv2/internal/core/runtime"
)

func TestSandboxGateShellDenied(t *testing.T) {
	manager := newManager(t, &policy.Bundle{
		BundleID:              "pb_shell",
		Name:                  "Shell Deny",
		Version:               "1",
		RulesetDigest:         "sha256:shell",
		DefaultVerdict:        policy.VerdictAllow,
		DefaultSandboxProfile: "local-default",
		Sandboxes: map[string]policy.SandboxProfile{
			"local-default": {Name: "local-default"},
		},
		Capabilities: []policy.CapabilityPolicy{
			{Name: "shell.exec", Verdict: policy.VerdictAllow, SandboxProfile: "local-default"},
		},
	})
	resp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: writeRuntimeContract(t, t.TempDir(), map[string]string{
		"requires_shell": "true",
	})})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if resp.Session.State != runtime.SessionStateFailedSafe {
		t.Fatalf("expected failed_safe, got %s", resp.Session.State)
	}
}
