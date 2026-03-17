package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/gov"
	"srosv2/internal/core/orch"
)

type Options struct {
	StoreDir     string
	Mode         string
	Gate         AdmissionGate
	Now          func() time.Time
	Orchestrator *orch.Orchestrator
	Governor     *gov.Engine
}

type Manager struct {
	store        *Store
	mode         string
	gate         AdmissionGate
	now          func() time.Time
	orchestrator *orch.Orchestrator
	governor     *gov.Engine
}

func NewManager(opts Options) (*Manager, error) {
	store, err := NewStore(opts.StoreDir)
	if err != nil {
		return nil, err
	}

	now := opts.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}

	mode := opts.Mode
	if strings.TrimSpace(mode) == "" {
		mode = "local_cli"
	}

	return &Manager{
		store:        store,
		mode:         mode,
		gate:         opts.Gate,
		now:          now,
		orchestrator: opts.Orchestrator,
		governor:     opts.Governor,
	}, nil
}

func (m *Manager) Run(ctx context.Context, req RunRequest) (RuntimeResponse, error) {
	contractPath := strings.TrimSpace(req.ContractPath)
	if contractPath == "" {
		return RuntimeResponse{}, fmt.Errorf("contract path is required")
	}

	contract, err := loadRunContract(contractPath)
	if err != nil {
		return RuntimeResponse{}, err
	}

	decision := AdmissionDecision{InitialState: SessionStateApproved, AutoStart: true, Reason: "runtime admission allow"}
	if m.gate != nil {
		decision, err = m.gate.Admit(ctx, AdmissionRequest{Contract: contract, ContractPath: contractPath})
		if err != nil {
			return RuntimeResponse{}, err
		}
	}

	now := m.now().UTC()
	session := NewSession(contract, contractPath, now)
	session.TopologyBinding = decision.TopologyBinding

	var plan *orch.Plan
	if m.orchestrator != nil {
		hydrated, err := m.orchestrator.Hydrate(session.SessionID, contract, decision.TopologyBinding)
		if err != nil {
			return RuntimeResponse{}, err
		}
		plan = &hydrated
		session.PlanPath = filepath.Join(m.store.Root(), "orch", session.SessionID+"_plan.json")
	}

	if plan != nil && m.governor != nil {
		executed, err := m.orchestrator.Execute(ctx, *plan, func(ctx context.Context, unit orch.WorkUnit) (orch.Decision, error) {
			result, err := m.governor.Evaluate(ctx, gov.Request{
				RunID:      contract.RunID,
				TraceID:    contract.TraceID,
				RiskClass:  contract.RiskClass,
				Capability: unit.Capability,
			})
			if err != nil {
				return orch.Decision{}, err
			}
			return orch.Decision{
				Verdict:        string(result.Decision.Verdict),
				Reason:         result.Decision.Reason,
				SandboxProfile: result.Decision.SandboxProfile,
			}, nil
		})
		if err != nil {
			return RuntimeResponse{}, err
		}
		session.LastDecision = executed.Decision.Verdict
		switch executed.Decision.Verdict {
		case "ask":
			session.ApprovalPath = executed.Route.ApprovalPath
			if err := Transition(&session, SessionStateWaitingForInput, executed.Decision.Reason, now); err != nil {
				return RuntimeResponse{}, err
			}
			if err := m.store.SaveApproval(ApprovalCheckpoint{
				SessionID:   session.SessionID,
				Reason:      executed.Decision.Reason,
				Approved:    false,
				RequestedAt: now.Format(time.RFC3339),
			}); err != nil {
				return RuntimeResponse{}, err
			}
		case "deny":
			if err := Transition(&session, SessionStateFailedSafe, executed.Decision.Reason, now); err != nil {
				return RuntimeResponse{}, err
			}
		}
	}

	if session.State == SessionStateWaitingForInput || session.State == SessionStateFailedSafe {
		if err := m.store.SaveSession(session); err != nil {
			return RuntimeResponse{}, err
		}
		summary := "runtime session governed"
		if session.State == SessionStateWaitingForInput {
			summary = "runtime waiting for operator input"
		}
		return RuntimeResponse{
			Accepted:      true,
			Summary:       summary,
			Session:       RefFromSession(session),
			ApprovalPath:  session.ApprovalPath,
			RuntimeRecord: m.store.Root(),
			Decision:      session.LastDecision,
			Plan:          plan,
		}, nil
	}

	switch decision.InitialState {
	case SessionStateApproved:
		if err := Transition(&session, SessionStateApproved, decision.Reason, now); err != nil {
			return RuntimeResponse{}, err
		}
		if decision.AutoStart {
			if err := Transition(&session, SessionStateRunning, "runtime session started", now); err != nil {
				return RuntimeResponse{}, err
			}
		}
	case SessionStateWaitingForInput:
		if err := Transition(&session, SessionStateWaitingForInput, decision.Reason, now); err != nil {
			return RuntimeResponse{}, err
		}
		approval := ApprovalCheckpoint{
			SessionID:   session.SessionID,
			Reason:      decision.WaitingApprovalHint,
			Approved:    false,
			RequestedAt: now.Format(time.RFC3339),
		}
		if err := m.store.SaveApproval(approval); err != nil {
			return RuntimeResponse{}, err
		}
		session.ApprovalPath = m.store.ApprovalPath(session.SessionID)
	case SessionStateFailedSafe:
		if err := Transition(&session, SessionStateFailedSafe, decision.Reason, now); err != nil {
			return RuntimeResponse{}, err
		}
	default:
		return RuntimeResponse{}, fmt.Errorf("unsupported admission initial state %s", decision.InitialState)
	}

	if err := m.store.SaveSession(session); err != nil {
		return RuntimeResponse{}, err
	}

	summary := "runtime session admitted"
	if session.State == SessionStateWaitingForInput {
		summary = "runtime waiting for operator input"
	}

	return RuntimeResponse{
		Accepted:      true,
		Summary:       summary,
		Session:       RefFromSession(session),
		ApprovalPath:  session.ApprovalPath,
		RuntimeRecord: m.store.Root(),
		Decision:      session.LastDecision,
		Plan:          plan,
	}, nil
}

func (m *Manager) Plan(context.Context, RunRequest) (RuntimeResponse, error) {
	return RuntimeResponse{Accepted: true, Summary: "runtime plan acknowledged; orchestration deferred to W06"}, nil
}

func (m *Manager) Resume(_ context.Context, req ResumeRequest) (RuntimeResponse, error) {
	session, err := m.getSession(req.SessionID)
	if err != nil {
		return RuntimeResponse{}, err
	}
	now := m.now().UTC()

	switch session.State {
	case SessionStatePaused:
		if err := Transition(&session, SessionStateRunning, "operator resume", now); err != nil {
			return RuntimeResponse{}, err
		}
	case SessionStateCheckpointed:
		if err := Transition(&session, SessionStateRunning, "resume from checkpoint", now); err != nil {
			return RuntimeResponse{}, err
		}
	case SessionStateWaitingForInput:
		approved, err := m.resolveApproval(session, req.ApprovalFile)
		if err != nil {
			return RuntimeResponse{}, err
		}
		if !approved {
			return RuntimeResponse{}, fmt.Errorf("resume requires explicit operator approval while waiting_for_input")
		}
		if err := Transition(&session, SessionStateApproved, "operator approval acknowledged", now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := Transition(&session, SessionStateRunning, "resume after approval", now); err != nil {
			return RuntimeResponse{}, err
		}
	default:
		return RuntimeResponse{}, fmt.Errorf("cannot resume session in state %s", session.State)
	}

	if err := m.store.SaveSession(session); err != nil {
		return RuntimeResponse{}, err
	}

	return RuntimeResponse{Accepted: true, Summary: "runtime resumed", Session: RefFromSession(session)}, nil
}

func (m *Manager) Pause(_ context.Context, req PauseRequest) (RuntimeResponse, error) {
	session, err := m.getSession(req.SessionID)
	if err != nil {
		return RuntimeResponse{}, err
	}
	now := m.now().UTC()
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "operator pause"
	}

	if err := Transition(&session, SessionStatePaused, reason, now); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.store.SaveSession(session); err != nil {
		return RuntimeResponse{}, err
	}

	return RuntimeResponse{Accepted: true, Summary: "runtime paused", Session: RefFromSession(session)}, nil
}

func (m *Manager) Checkpoint(_ context.Context, req CheckpointRequest) (RuntimeResponse, error) {
	session, err := m.getSession(req.SessionID)
	if err != nil {
		return RuntimeResponse{}, err
	}
	now := m.now().UTC()
	cp, err := NewCheckpoint(session, req.Stage, now)
	if err != nil {
		return RuntimeResponse{}, err
	}

	if err := Transition(&session, SessionStateCheckpointed, "checkpoint created", now); err != nil {
		return RuntimeResponse{}, err
	}
	session.LatestCheckpointID = string(cp.Record.CheckpointID)

	if err := m.store.SaveCheckpoint(cp); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.store.SaveSession(session); err != nil {
		return RuntimeResponse{}, err
	}

	return RuntimeResponse{
		Accepted:     true,
		Summary:      "checkpoint created",
		Session:      RefFromSession(session),
		CheckpointID: string(cp.Record.CheckpointID),
	}, nil
}

func (m *Manager) Rollback(_ context.Context, req RollbackRequest) (RuntimeResponse, error) {
	session, err := m.getSession(req.SessionID)
	if err != nil {
		return RuntimeResponse{}, err
	}

	checkpointID := strings.TrimSpace(req.CheckpointID)
	if checkpointID == "" {
		checkpointID = session.LatestCheckpointID
	}
	if strings.TrimSpace(checkpointID) == "" {
		return RuntimeResponse{}, fmt.Errorf("rollback requires checkpoint id")
	}
	if _, err := m.store.LoadCheckpoint(checkpointID); err != nil {
		return RuntimeResponse{}, fmt.Errorf("load checkpoint %s: %w", checkpointID, err)
	}

	now := m.now().UTC()
	rb, err := NewRollback(session, checkpointID, req.Reason, now)
	if err != nil {
		return RuntimeResponse{}, err
	}

	if err := Transition(&session, SessionStateRolledBack, "rollback to checkpoint "+checkpointID, now); err != nil {
		return RuntimeResponse{}, err
	}
	session.LatestRollbackID = string(rb.Record.RollbackID)

	if err := m.store.SaveRollback(rb); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.store.SaveSession(session); err != nil {
		return RuntimeResponse{}, err
	}

	return RuntimeResponse{
		Accepted:   true,
		Summary:    "runtime rolled back",
		Session:    RefFromSession(session),
		RollbackID: string(rb.Record.RollbackID),
	}, nil
}

func (m *Manager) Status(_ context.Context, req StatusRequest) (StatusSnapshot, error) {
	sessionID := strings.TrimSpace(req.SessionID)
	if req.Latest || sessionID == "" {
		s, err := m.store.LatestSession()
		if err != nil {
			return StatusSnapshot{Mode: m.mode, Summary: "no runtime session found"}, nil
		}
		return snapshotFromSession(m.mode, s), nil
	}

	session, err := m.store.LoadSession(sessionID)
	if err != nil {
		return StatusSnapshot{}, err
	}
	return snapshotFromSession(m.mode, session), nil
}

func snapshotFromSession(mode string, session RuntimeSession) StatusSnapshot {
	return StatusSnapshot{
		Mode:               mode,
		Session:            ptrSessionRef(RefFromSession(session)),
		Summary:            session.Reason,
		LatestCheckpointID: session.LatestCheckpointID,
		LatestRollbackID:   session.LatestRollbackID,
		WaitingApproval:    session.ApprovalPath,
		PlanPath:           session.PlanPath,
		LastDecision:       session.LastDecision,
	}
}

func ptrSessionRef(v SessionRef) *SessionRef { return &v }

func (m *Manager) getSession(sessionID string) (RuntimeSession, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return m.store.LatestSession()
	}
	return m.store.LoadSession(sessionID)
}

func (m *Manager) resolveApproval(session RuntimeSession, approvalFile string) (bool, error) {
	if strings.TrimSpace(approvalFile) != "" {
		data, err := os.ReadFile(approvalFile)
		if err != nil {
			return false, fmt.Errorf("read approval file: %w", err)
		}
		var payload struct {
			Approved bool `json:"approved"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return false, fmt.Errorf("decode approval file: %w", err)
		}
		return payload.Approved, nil
	}

	approval, err := m.store.LoadApproval(session.SessionID)
	if err != nil {
		return false, nil
	}
	return approval.Approved, nil
}

func (m *Manager) ToolsList(context.Context) (map[string]any, error) {
	bundle := gov.DefaultBundle()
	if m.governor != nil {
		bundle = m.governor.Bundle()
	}
	capabilities := make([]map[string]any, 0, len(bundle.Capabilities))
	for _, item := range bundle.Capabilities {
		capabilities = append(capabilities, map[string]any{
			"name":            item.Name,
			"verdict":         item.Verdict,
			"sandbox_profile": item.SandboxProfile,
			"deferred_to":     "W09",
		})
	}
	return map[string]any{
		"bundle_id":      bundle.BundleID,
		"capabilities":   capabilities,
		"execution_mode": "governed_semantics_only",
	}, nil
}

func (m *Manager) ToolsShow(_ context.Context, name string) (map[string]any, error) {
	bundle := gov.DefaultBundle()
	if m.governor != nil {
		bundle = m.governor.Bundle()
	}
	for _, item := range bundle.Capabilities {
		if item.Name == strings.TrimSpace(name) {
			return map[string]any{
				"name":               item.Name,
				"verdict":            item.Verdict,
				"sandbox_profile":    item.SandboxProfile,
				"allowed_boundaries": item.AllowedBoundaries,
				"deferred_to":        "W09",
			}, nil
		}
	}
	return nil, fmt.Errorf("tool capability %q not found", name)
}

func (m *Manager) ToolsValidate(_ context.Context, path string) (map[string]any, error) {
	bundle, err := gov.LoadBundle(strings.TrimSpace(path))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"valid":        true,
		"bundle_id":    bundle.BundleID,
		"name":         bundle.Name,
		"version":      bundle.Version,
		"capabilities": len(bundle.Capabilities),
		"sandboxes":    len(bundle.Sandboxes),
	}, nil
}

func (m *Manager) ToolsRegister(_ context.Context, path string) (map[string]any, error) {
	bundle, err := gov.LoadBundle(strings.TrimSpace(path))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"registered": false,
		"bundle_id":  bundle.BundleID,
		"summary":    "registration remains deferred to W09; bundle validated through GOV",
	}, nil
}

func (m *Manager) ConnectorsList(context.Context) (map[string]any, error) {
	return map[string]any{
		"connectors": []map[string]any{
			{"name": "local-governed-connector", "verdict": "ask", "deferred_to": "W09"},
		},
		"execution_mode": "governed_semantics_only",
	}, nil
}

func (m *Manager) MCPIngest(_ context.Context, path string) (map[string]any, error) {
	return map[string]any{
		"accepted":       false,
		"file":           strings.TrimSpace(path),
		"summary":        "MCP transport remains deferred to W09; GOV can govern the capability boundary only",
		"execution_mode": "governed_semantics_only",
	}, nil
}

func loadRunContract(path string) (runcontract.RunContract, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return runcontract.RunContract{}, fmt.Errorf("read run contract file %q: %w", path, err)
	}

	var contract runcontract.RunContract
	if err := json.Unmarshal(data, &contract); err != nil {
		return runcontract.RunContract{}, fmt.Errorf("decode run contract file %q: %w", path, err)
	}
	if errs := runcontract.Validate(contract); len(errs) > 0 {
		return runcontract.RunContract{}, fmt.Errorf("invalid run contract: %v", errs[0])
	}

	return contract, nil
}
