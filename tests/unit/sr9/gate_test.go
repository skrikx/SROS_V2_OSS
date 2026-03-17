package sr9_test

import (
	"context"
	"testing"

	"srosv2/internal/core/runtime"
	"srosv2/internal/core/sr9"
)

func TestGateAdmitAllow(t *testing.T) {
	gate := sr9.NewGate(sr9.Options{})
	decision, err := gate.Admit(context.Background(), runtime.AdmissionRequest{Contract: validContract()})
	if err != nil {
		t.Fatalf("admit: %v", err)
	}
	if decision.InitialState != runtime.SessionStateApproved {
		t.Fatalf("expected approved state, got %s", decision.InitialState)
	}
	if !decision.AutoStart {
		t.Fatal("allow decision should autostart")
	}
}

func TestGateAdmitAsk(t *testing.T) {
	contract := validContract()
	contract.Metadata["approval_mode"] = "ask"

	gate := sr9.NewGate(sr9.Options{})
	decision, err := gate.Admit(context.Background(), runtime.AdmissionRequest{Contract: contract})
	if err != nil {
		t.Fatalf("admit: %v", err)
	}
	if decision.InitialState != runtime.SessionStateWaitingForInput {
		t.Fatalf("expected waiting_for_input state, got %s", decision.InitialState)
	}
	if !decision.RequireOperatorAck {
		t.Fatal("ask decision should require operator ack")
	}
}
