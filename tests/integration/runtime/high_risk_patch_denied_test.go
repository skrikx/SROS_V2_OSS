package runtime_integration_test

import (
	"context"
	"testing"

	"srosv2/contracts/policy"
	"srosv2/internal/core/runtime"
)

func TestHighRiskPatchDenied(t *testing.T) {
	manager := newManager(t, &policy.Bundle{
		BundleID:              "pb_patch",
		Name:                  "Patch Deny",
		Version:               "1",
		RulesetDigest:         "sha256:patch",
		DefaultVerdict:        policy.VerdictAllow,
		DefaultSandboxProfile: "patch-safe",
		Sandboxes: map[string]policy.SandboxProfile{
			"patch-safe": {Name: "patch-safe", AllowPatch: true},
		},
	})
	resp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: writeRuntimeContract(t, t.TempDir(), map[string]string{
		"requires_patch": "true",
	})})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if resp.Session.State != runtime.SessionStateFailedSafe {
		t.Fatalf("expected failed_safe, got %s", resp.Session.State)
	}
}
