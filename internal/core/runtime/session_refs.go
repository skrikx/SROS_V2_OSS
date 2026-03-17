package runtime

import "time"

type SessionState string

const (
	SessionStatePlanned         SessionState = "planned"
	SessionStateApproved        SessionState = "approved"
	SessionStateRunning         SessionState = "running"
	SessionStateWaitingForInput SessionState = "waiting_for_input"
	SessionStatePaused          SessionState = "paused"
	SessionStateCheckpointed    SessionState = "checkpointed"
	SessionStateFailedSafe      SessionState = "failed_safe"
	SessionStateRolledBack      SessionState = "rolled_back"
	SessionStateCompleted       SessionState = "completed"
)

type SessionRef struct {
	RunID     string       `json:"run_id"`
	SessionID string       `json:"session_id"`
	State     SessionState `json:"state"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type StatusSnapshot struct {
	Mode               string      `json:"mode"`
	Session            *SessionRef `json:"session,omitempty"`
	Summary            string      `json:"summary"`
	Boundaries         []string    `json:"boundaries,omitempty"`
	LatestCheckpointID string      `json:"latest_checkpoint_id,omitempty"`
	LatestRollbackID   string      `json:"latest_rollback_id,omitempty"`
	WaitingApproval    string      `json:"waiting_approval,omitempty"`
}
