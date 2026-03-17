package runtime

import (
	"context"

	"srosv2/contracts/runcontract"
)

type RunRequest struct {
	ContractPath string `json:"contract_path"`
}

type ResumeRequest struct {
	SessionID    string `json:"session_id,omitempty"`
	ApprovalFile string `json:"approval_file,omitempty"`
}

type PauseRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type CheckpointRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Stage     string `json:"stage,omitempty"`
}

type RollbackRequest struct {
	SessionID    string `json:"session_id,omitempty"`
	CheckpointID string `json:"checkpoint_id,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

type RuntimeResponse struct {
	Accepted      bool       `json:"accepted"`
	Summary       string     `json:"summary"`
	Session       SessionRef `json:"session,omitempty"`
	CheckpointID  string     `json:"checkpoint_id,omitempty"`
	RollbackID    string     `json:"rollback_id,omitempty"`
	ApprovalPath  string     `json:"approval_path,omitempty"`
	RuntimeRecord string     `json:"runtime_record,omitempty"`
}

type StatusRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Latest    bool   `json:"latest,omitempty"`
}

type AdmissionRequest struct {
	Contract     runcontract.RunContract `json:"contract"`
	ContractPath string                  `json:"contract_path"`
}

type AdmissionDecision struct {
	InitialState        SessionState `json:"initial_state"`
	Reason              string       `json:"reason"`
	AutoStart           bool         `json:"auto_start"`
	RequireOperatorAck  bool         `json:"require_operator_ack"`
	WaitingApprovalHint string       `json:"waiting_approval_hint,omitempty"`
	TopologyBinding     string       `json:"topology_binding,omitempty"`
}

type AdmissionGate interface {
	Admit(context.Context, AdmissionRequest) (AdmissionDecision, error)
}

type Runtime interface {
	Run(context.Context, RunRequest) (RuntimeResponse, error)
	Plan(context.Context, RunRequest) (RuntimeResponse, error)
	Resume(context.Context, ResumeRequest) (RuntimeResponse, error)
	Pause(context.Context, PauseRequest) (RuntimeResponse, error)
	Checkpoint(context.Context, CheckpointRequest) (RuntimeResponse, error)
	Rollback(context.Context, RollbackRequest) (RuntimeResponse, error)
	Status(context.Context, StatusRequest) (StatusSnapshot, error)
}

type Inspector interface {
	Trace(context.Context, string) (map[string]any, error)
	Receipts(context.Context, string) (map[string]any, error)
	Memory(context.Context, string) (map[string]any, error)
	Mirror(context.Context, string) (map[string]any, error)
}

type Fabric interface {
	ToolsList(context.Context) (map[string]any, error)
	ToolsShow(context.Context, string) (map[string]any, error)
	ToolsValidate(context.Context, string) (map[string]any, error)
	ToolsRegister(context.Context, string) (map[string]any, error)
	ConnectorsList(context.Context) (map[string]any, error)
	MCPIngest(context.Context, string) (map[string]any, error)
}
