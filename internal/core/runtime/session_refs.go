package runtime

type SessionState string

const (
	SessionStateUnknown SessionState = "unknown"
	SessionStateReady   SessionState = "ready"
	SessionStateRunning SessionState = "running"
	SessionStatePaused  SessionState = "paused"
	SessionStateClosed  SessionState = "closed"
)

type SessionRef struct {
	RunID     string       `json:"run_id"`
	SessionID string       `json:"session_id"`
	State     SessionState `json:"state"`
}

type StatusSnapshot struct {
	Mode       string      `json:"mode"`
	Session    *SessionRef `json:"session,omitempty"`
	Summary    string      `json:"summary"`
	Boundaries []string    `json:"boundaries,omitempty"`
}
