package runtime_test

import (
	"strings"
	"testing"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/runtime"
)

func TestNewSessionInitializesCanonicalFields(t *testing.T) {
	contract := runcontract.RunContract{RunID: "run_session_test"}
	session := runtime.NewSession(contract, "contract.json", fixedNow)

	if session.State != runtime.SessionStatePlanned {
		t.Fatalf("expected planned state, got %s", session.State)
	}
	if session.RunID != "run_session_test" {
		t.Fatalf("unexpected run id %s", session.RunID)
	}
	if !strings.HasPrefix(session.SessionID, "sess_") {
		t.Fatalf("expected session id prefix sess_, got %s", session.SessionID)
	}
}

func TestRefFromSession(t *testing.T) {
	contract := runcontract.RunContract{RunID: "run_session_ref"}
	session := runtime.NewSession(contract, "contract.json", fixedNow)
	ref := runtime.RefFromSession(session)

	if ref.SessionID != session.SessionID {
		t.Fatalf("expected session id %s, got %s", session.SessionID, ref.SessionID)
	}
	if ref.State != session.State {
		t.Fatalf("expected state %s, got %s", session.State, ref.State)
	}
}
