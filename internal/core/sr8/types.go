package sr8

import (
	"time"

	"srosv2/contracts/evidence"
	"srosv2/contracts/runcontract"
)

type IntentSource string

const (
	IntentSourceInline IntentSource = "inline"
	IntentSourceFile   IntentSource = "file"
)

type DomainClass string

const (
	DomainGeneral  DomainClass = "general_ops"
	DomainFileTask DomainClass = "file_patch"
	DomainResearch DomainClass = "research"
)

type TopologyClass string

const (
	TopologyLocalSingle     TopologyClass = "local_single"
	TopologyLocalFilesystem TopologyClass = "local_filesystem"
	TopologyLocalResearch   TopologyClass = "local_research"
)

type CompileRequest struct {
	Intent       string `json:"intent,omitempty"`
	InputPath    string `json:"input_path,omitempty"`
	OperatorID   string `json:"operator_id,omitempty"`
	TenantID     string `json:"tenant_id,omitempty"`
	WorkspaceID  string `json:"workspace_id,omitempty"`
	ArtifactRoot string `json:"artifact_root,omitempty"`
	EmitSRXML    bool   `json:"emit_srxml,omitempty"`
}

type ParsedIntent struct {
	Source      IntentSource `json:"source"`
	RawIntent   string       `json:"raw_intent"`
	InputPath   string       `json:"input_path,omitempty"`
	OperatorID  string       `json:"operator_id"`
	TenantID    string       `json:"tenant_id"`
	WorkspaceID string       `json:"workspace_id"`
}

type NormalizedIntent struct {
	Source           IntentSource `json:"source"`
	RawIntent        string       `json:"raw_intent"`
	NormalizedText   string       `json:"normalized_text"`
	IntentSummary    string       `json:"intent_summary"`
	InputPath        string       `json:"input_path,omitempty"`
	OperatorID       string       `json:"operator_id"`
	TenantID         string       `json:"tenant_id"`
	WorkspaceID      string       `json:"workspace_id"`
	CompileRequestID string       `json:"compile_request_id"`
	RequestedAt      time.Time    `json:"requested_at"`
}

type Classification struct {
	Domain          DomainClass           `json:"domain"`
	Risk            runcontract.RiskClass `json:"risk"`
	ArtifactPosture string                `json:"artifact_posture"`
	Signals         []string              `json:"signals,omitempty"`
}

type TopologyDecision struct {
	Topology   TopologyClass          `json:"topology"`
	RouteClass runcontract.RouteClass `json:"route_class"`
	Reason     string                 `json:"reason"`
}

type CompileReceipt struct {
	Receipt          evidence.Receipt      `json:"receipt"`
	Bundle           evidence.Bundle       `json:"bundle"`
	CompileRequestID string                `json:"compile_request_id"`
	DomainClass      DomainClass           `json:"domain_class"`
	RiskClass        runcontract.RiskClass `json:"risk_class"`
	TopologyClass    TopologyClass         `json:"topology_class"`
	Status           string                `json:"status"`
	TraceLinkage     string                `json:"trace_linkage"`
	RuntimeAdmission string                `json:"runtime_admission"`
}

type CompileResult struct {
	Accepted       bool                            `json:"accepted"`
	Summary        string                          `json:"summary"`
	Normalized     NormalizedIntent                `json:"normalized"`
	Classification Classification                  `json:"classification"`
	Topology       TopologyDecision                `json:"topology"`
	RunContract    runcontract.RunContract         `json:"run_contract"`
	SRXML          string                          `json:"srxml,omitempty"`
	Receipt        CompileReceipt                  `json:"receipt"`
	Artifacts      []runcontract.ArtifactReference `json:"artifacts"`
	OutputDir      string                          `json:"output_dir"`
}
