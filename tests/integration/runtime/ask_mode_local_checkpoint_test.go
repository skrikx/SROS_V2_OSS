package runtime_integration_test

import (
	"context"
	"testing"

	"srosv2/contracts/policy"
	"srosv2/internal/core/runtime"
)

func TestAskModeCreatesLocalCheckpoint(t *testing.T) {
	manager := newManager(t, &policy.Bundle{
		BundleID:              "pb_ask",
		Name:                  "Ask",
		Version:               "1",
		RulesetDigest:         "sha256:ask",
		DefaultVerdict:        policy.VerdictAllow,
		DefaultSandboxProfile: "net-observe",
		Sandboxes: map[string]policy.SandboxProfile{
			"net-observe": {Name: "net-observe", AllowExternalNet: true},
		},
		Capabilities: []policy.CapabilityPolicy{
			{Name: "connector.invoke", Verdict: policy.VerdictAsk, SandboxProfile: "net-observe", RequireApproval: true},
		},
	})
	resp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: writeRuntimeContract(t, t.TempDir(), map[string]string{
		"requires_connector": "true",
	})})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if resp.Session.State != runtime.SessionStateWaitingForInput {
		t.Fatalf("expected waiting_for_input, got %s", resp.Session.State)
	}
	if resp.ApprovalPath == "" {
		t.Fatal("expected approval artifact path")
	}
}
