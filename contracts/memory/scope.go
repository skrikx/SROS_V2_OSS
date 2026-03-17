package memory

import "srosv2/internal/shared/ids"

type Scope string

const (
	ScopeSession   Scope = "session"
	ScopeWorkspace Scope = "workspace"
	ScopeRun       Scope = "run"
	ScopeGlobal    Scope = "global"
)

type ScopeBinding struct {
	Scope       Scope           `json:"scope"`
	TenantID    ids.TenantID    `json:"tenant_id"`
	WorkspaceID ids.WorkspaceID `json:"workspace_id"`
	RunID       ids.RunID       `json:"run_id,omitempty"`
	SessionID   ids.SessionID   `json:"session_id,omitempty"`
}
