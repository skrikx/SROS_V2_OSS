package runtime_test

import (
	"testing"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/runtime"
)

func TestTransitionUpdatesSessionHistory(t *testing.T) {
	session := runtime.NewSession(runcontract.RunContract{RunID: "run_1"}, "contract.json", fixedNow)
	if session.State != runtime.SessionStatePlanned {
		t.Fatalf("expected planned initial state, got %s", session.State)
	}

	if err := runtime.Transition(&session, runtime.SessionStateApproved, "approved for runtime", fixedNow); err != nil {
		t.Fatalf("transition to approved: %v", err)
	}

	if session.State != runtime.SessionStateApproved {
		t.Fatalf("expected approved state, got %s", session.State)
	}
	if len(session.History) != 1 {
		t.Fatalf("expected one lifecycle event, got %d", len(session.History))
	}
	if session.History[0].From != runtime.SessionStatePlanned || session.History[0].To != runtime.SessionStateApproved {
		t.Fatalf("unexpected transition record: %+v", session.History[0])
	}
}

func TestTransitionRejectsIllegalMove(t *testing.T) {
	session := runtime.NewSession(runcontract.RunContract{RunID: "run_2"}, "contract.json", fixedNow)
	if err := runtime.Transition(&session, runtime.SessionStateRolledBack, "illegal", fixedNow); err == nil {
		t.Fatal("expected illegal transition error")
	}
	if session.State != runtime.SessionStatePlanned {
		t.Fatalf("state should remain planned after illegal transition, got %s", session.State)
	}
}
