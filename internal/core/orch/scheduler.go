package orch

import (
	"fmt"
	"strings"

	"srosv2/contracts/runcontract"
)

type WorkUnit struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Capability  string   `json:"capability"`
	Privileged  bool     `json:"privileged"`
	Checkpoint  bool     `json:"checkpoint"`
	DependsOn   []string `json:"depends_on,omitempty"`
	Description string   `json:"description"`
}

type Plan struct {
	SessionID    string           `json:"session_id"`
	RunID        string           `json:"run_id"`
	RiskClass    string           `json:"risk_class"`
	Concurrency  ConcurrencyRules `json:"concurrency"`
	WorkUnits    []WorkUnit       `json:"work_units"`
	Artifacts    []string         `json:"artifacts,omitempty"`
	TopologyHint string           `json:"topology_hint,omitempty"`
}

type Scheduler struct{}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Build(sessionID string, contract runcontract.RunContract, topologyHint string) (Plan, error) {
	if strings.TrimSpace(sessionID) == "" {
		return Plan{}, fmt.Errorf("session id is required")
	}

	rules := RulesForRisk(contract.RiskClass)
	workUnits := []WorkUnit{
		{
			ID:          "wu-001",
			Name:        "hydrate-runtime-session",
			Capability:  "runtime.session.prepare",
			Privileged:  false,
			Description: "hydrate runtime session into an executable orchestration plan",
		},
	}

	lastID := "wu-001"
	add := func(name, capability, description string, privileged, checkpoint bool) {
		id := fmt.Sprintf("wu-%03d", len(workUnits)+1)
		workUnits = append(workUnits, WorkUnit{
			ID:          id,
			Name:        name,
			Capability:  capability,
			Privileged:  privileged,
			Checkpoint:  checkpoint,
			DependsOn:   []string{lastID},
			Description: description,
		})
		lastID = id
	}

	if strings.EqualFold(contract.Metadata["requires_shell"], "true") {
		add("govern-shell-capability", "shell.exec", "govern shell execution before any runtime side effect", true, false)
	}
	if strings.EqualFold(contract.Metadata["requires_patch"], "true") {
		add("govern-patch-capability", "patch.apply", "govern patch capability before local filesystem mutation", true, false)
	}
	if strings.EqualFold(contract.Metadata["requires_tool_validation"], "true") {
		add("govern-tool-validation", "tool.validate", "govern tool manifest validation boundary", true, false)
	}
	if strings.EqualFold(contract.Metadata["requires_connector"], "true") {
		add("govern-connector-boundary", "connector.invoke", "govern connector boundary without executing it", true, false)
	}
	if strings.EqualFold(contract.Metadata["requires_mcp"], "true") {
		add("govern-mcp-boundary", "mcp.ingest", "govern MCP boundary without transport execution", true, false)
	}

	for idx, ref := range contract.CheckpointRefs {
		add(
			fmt.Sprintf("checkpoint-%d", idx+1),
			"runtime.checkpoint.route",
			"route a local operator checkpoint for "+string(ref.CheckpointID),
			false,
			true,
		)
	}

	add("coordinate-artifact-emission", "runtime.artifacts.coordinate", "coordinate artifact emission and receipt hooks", false, false)

	return Plan{
		SessionID:    sessionID,
		RunID:        string(contract.RunID),
		RiskClass:    string(contract.RiskClass),
		Concurrency:  rules,
		WorkUnits:    workUnits,
		Artifacts:    []string{"queue_plan.json", "governance_decisions.json"},
		TopologyHint: topologyHint,
	}, nil
}
