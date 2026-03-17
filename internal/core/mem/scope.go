package mem

import (
	"fmt"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/shared/ids"
)

type ScopeBinding struct {
	Scope       cmemory.Scope   `json:"scope"`
	TenantID    ids.TenantID    `json:"tenant_id"`
	WorkspaceID ids.WorkspaceID `json:"workspace_id"`
	RunID       ids.RunID       `json:"run_id,omitempty"`
	SessionID   ids.SessionID   `json:"session_id,omitempty"`
}

func (s ScopeBinding) Validate() error {
	if s.Scope == "" {
		return fmt.Errorf("scope is required")
	}
	if s.TenantID == "" {
		return fmt.Errorf("tenant id is required")
	}
	if s.WorkspaceID == "" {
		return fmt.Errorf("workspace id is required")
	}
	return nil
}

func (s ScopeBinding) Contract() cmemory.ScopeBinding {
	return cmemory.ScopeBinding{
		Scope:       s.Scope,
		TenantID:    s.TenantID,
		WorkspaceID: s.WorkspaceID,
		RunID:       s.RunID,
		SessionID:   s.SessionID,
	}
}
