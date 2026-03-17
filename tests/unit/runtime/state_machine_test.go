package runtime_test

import (
	"testing"

	"srosv2/internal/core/runtime"
)

func TestMinimumLegalTransitions(t *testing.T) {
	cases := []struct {
		from runtime.SessionState
		to   runtime.SessionState
	}{
		{runtime.SessionStatePlanned, runtime.SessionStateApproved},
		{runtime.SessionStatePlanned, runtime.SessionStateWaitingForInput},
		{runtime.SessionStatePlanned, runtime.SessionStateFailedSafe},
		{runtime.SessionStateApproved, runtime.SessionStateRunning},
		{runtime.SessionStateApproved, runtime.SessionStateWaitingForInput},
		{runtime.SessionStateApproved, runtime.SessionStateFailedSafe},
		{runtime.SessionStateRunning, runtime.SessionStateWaitingForInput},
		{runtime.SessionStateRunning, runtime.SessionStatePaused},
		{runtime.SessionStateRunning, runtime.SessionStateCheckpointed},
		{runtime.SessionStateRunning, runtime.SessionStateCompleted},
		{runtime.SessionStateRunning, runtime.SessionStateFailedSafe},
		{runtime.SessionStateWaitingForInput, runtime.SessionStateApproved},
		{runtime.SessionStateWaitingForInput, runtime.SessionStateRunning},
		{runtime.SessionStateWaitingForInput, runtime.SessionStatePaused},
		{runtime.SessionStateWaitingForInput, runtime.SessionStateFailedSafe},
		{runtime.SessionStatePaused, runtime.SessionStateRunning},
		{runtime.SessionStatePaused, runtime.SessionStateCheckpointed},
		{runtime.SessionStatePaused, runtime.SessionStateRolledBack},
		{runtime.SessionStatePaused, runtime.SessionStateFailedSafe},
		{runtime.SessionStateCheckpointed, runtime.SessionStateRunning},
		{runtime.SessionStateCheckpointed, runtime.SessionStatePaused},
		{runtime.SessionStateCheckpointed, runtime.SessionStateRolledBack},
		{runtime.SessionStateCheckpointed, runtime.SessionStateFailedSafe},
	}

	for _, tc := range cases {
		if !runtime.CanTransition(tc.from, tc.to) {
			t.Fatalf("expected legal transition %s -> %s", tc.from, tc.to)
		}
	}
}

func TestIllegalTransitionRejected(t *testing.T) {
	if runtime.CanTransition(runtime.SessionStateCompleted, runtime.SessionStateRunning) {
		t.Fatal("completed -> running must be illegal")
	}
	if err := runtime.EnsureTransition(runtime.SessionStateCompleted, runtime.SessionStateRunning); err == nil {
		t.Fatal("expected illegal transition error")
	}
}

func TestTerminalStates(t *testing.T) {
	if !runtime.IsTerminal(runtime.SessionStateFailedSafe) {
		t.Fatal("failed_safe should be terminal")
	}
	if !runtime.IsTerminal(runtime.SessionStateRolledBack) {
		t.Fatal("rolled_back should be terminal")
	}
	if !runtime.IsTerminal(runtime.SessionStateCompleted) {
		t.Fatal("completed should be terminal")
	}
	if runtime.IsTerminal(runtime.SessionStateRunning) {
		t.Fatal("running should not be terminal")
	}
}
