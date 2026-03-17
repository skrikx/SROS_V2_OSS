package runtime_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"srosv2/internal/core/runtime"
)

func TestManagerRunCreatesRuntimeSession(t *testing.T) {
	dir := t.TempDir()
	contractPath := writeRunContractFile(t, dir, nil)

	manager, err := runtime.NewManager(runtime.Options{
		StoreDir: dir,
		Mode:     "local_cli",
		Now:      func() time.Time { return fixedNow },
		Gate: stubGate{decision: runtime.AdmissionDecision{
			InitialState: runtime.SessionStateApproved,
			Reason:       "allow",
			AutoStart:    true,
		}},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	resp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: contractPath})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !resp.Accepted {
		t.Fatal("expected accepted response")
	}
	if resp.Session.State != runtime.SessionStateRunning {
		t.Fatalf("expected running state, got %s", resp.Session.State)
	}

	status, err := manager.Status(context.Background(), runtime.StatusRequest{Latest: true})
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if status.Session == nil {
		t.Fatal("expected latest session in status")
	}
	if status.Session.State != runtime.SessionStateRunning {
		t.Fatalf("expected running status state, got %s", status.Session.State)
	}
}

func TestManagerRunAskModeTransitionsWaitingForInput(t *testing.T) {
	dir := t.TempDir()
	contractPath := writeRunContractFile(t, dir, nil)

	manager, err := runtime.NewManager(runtime.Options{
		StoreDir: dir,
		Now:      func() time.Time { return fixedNow },
		Gate: stubGate{decision: runtime.AdmissionDecision{
			InitialState:        runtime.SessionStateWaitingForInput,
			Reason:              "ask",
			AutoStart:           false,
			RequireOperatorAck:  true,
			WaitingApprovalHint: "operator approval required",
		}},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	resp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: contractPath})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if resp.Session.State != runtime.SessionStateWaitingForInput {
		t.Fatalf("expected waiting_for_input, got %s", resp.Session.State)
	}
	if resp.ApprovalPath == "" {
		t.Fatal("expected approval path for ask-mode")
	}
}

func TestManagerResumeFromAskRequiresApproval(t *testing.T) {
	dir := t.TempDir()
	contractPath := writeRunContractFile(t, dir, nil)

	manager, err := runtime.NewManager(runtime.Options{
		StoreDir: dir,
		Now:      func() time.Time { return fixedNow },
		Gate: stubGate{decision: runtime.AdmissionDecision{
			InitialState:        runtime.SessionStateWaitingForInput,
			Reason:              "ask",
			AutoStart:           false,
			RequireOperatorAck:  true,
			WaitingApprovalHint: "operator approval required",
		}},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	runResp, err := manager.Run(context.Background(), runtime.RunRequest{ContractPath: contractPath})
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	if _, err := manager.Resume(context.Background(), runtime.ResumeRequest{SessionID: runResp.Session.SessionID}); err == nil {
		t.Fatal("expected resume to fail without approval")
	}

	approvalPath := filepath.Join(dir, "manual_approval.json")
	if err := os.WriteFile(approvalPath, []byte("{\"approved\":true}\n"), 0o644); err != nil {
		t.Fatalf("write approval file: %v", err)
	}

	resumeResp, err := manager.Resume(context.Background(), runtime.ResumeRequest{
		SessionID:    runResp.Session.SessionID,
		ApprovalFile: approvalPath,
	})
	if err != nil {
		t.Fatalf("resume with approval: %v", err)
	}
	if resumeResp.Session.State != runtime.SessionStateRunning {
		t.Fatalf("expected running state after approved resume, got %s", resumeResp.Session.State)
	}
}
