package runtime

import (
	"context"

	"srosv2/contracts/runcontract"
)

type CompileRequest struct {
	Intent    string `json:"intent"`
	InputPath string `json:"input_path,omitempty"`
}

type CompileResponse struct {
	Accepted    bool                     `json:"accepted"`
	Summary     string                   `json:"summary"`
	RunContract *runcontract.RunContract `json:"run_contract,omitempty"`
}

type RunRequest struct {
	RunID string `json:"run_id,omitempty"`
	Plan  string `json:"plan,omitempty"`
}

type ResumeRequest struct {
	SessionID string `json:"session_id"`
}

type PauseRequest struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason,omitempty"`
}

type CheckpointRequest struct {
	SessionID string `json:"session_id"`
	Stage     string `json:"stage"`
}

type RollbackRequest struct {
	SessionID    string `json:"session_id"`
	CheckpointID string `json:"checkpoint_id"`
}

type RuntimeResponse struct {
	Accepted bool       `json:"accepted"`
	Summary  string     `json:"summary"`
	Session  SessionRef `json:"session,omitempty"`
}

type StatusRequest struct {
	SessionID string `json:"session_id,omitempty"`
}

type Compiler interface {
	Compile(context.Context, CompileRequest) (CompileResponse, error)
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
