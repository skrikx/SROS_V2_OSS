package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"srosv2/contracts/evidence"
	cmemory "srosv2/contracts/memory"
	"srosv2/contracts/runcontract"
	ctrace "srosv2/contracts/trace"
	"srosv2/internal/core/gov"
	"srosv2/internal/core/mem"
	"srosv2/internal/core/mirror"
	"srosv2/internal/core/orch"
	coreprov "srosv2/internal/core/provenance"
	coretrace "srosv2/internal/core/trace"
	"srosv2/internal/shared/ids"
)

type Options struct {
	StoreDir     string
	Mode         string
	Gate         AdmissionGate
	Now          func() time.Time
	Orchestrator *orch.Orchestrator
	Governor     *gov.Engine
	Memory       *mem.Store
	Mirror       *mirror.Engine
	Trace        *coretrace.Service
	Provenance   *coreprov.Service
}

type Manager struct {
	store        *Store
	mode         string
	gate         AdmissionGate
	now          func() time.Time
	orchestrator *orch.Orchestrator
	governor     *gov.Engine
	memory       *mem.Store
	mirror       *mirror.Engine
	trace        *coretrace.Service
	provenance   *coreprov.Service
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
		memory:       opts.Memory,
		mirror:       opts.Mirror,
		trace:        opts.Trace,
		provenance:   opts.Provenance,
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

	if err := m.emitTrace(contract.RunID, contract.TraceID, ids.SpanID(""), coretrace.EventRunStarted, map[string]any{
		"session_id":    session.SessionID,
		"contract_path": contractPath,
		"risk_class":    contract.RiskClass,
	}); err != nil {
		return RuntimeResponse{}, err
	}

	if err := m.recordMemoryMutation(&session, "runtime.session", "created", "runtime session created"); err != nil {
		return RuntimeResponse{}, err
	}

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
			if err := m.emitTrace(contract.RunID, contract.TraceID, ids.SpanID(""), coretrace.EventPolicyDecision, map[string]any{"verdict": "deny", "reason": executed.Decision.Reason}); err != nil {
				return RuntimeResponse{}, err
			}
		}
	}

	if session.State == SessionStateWaitingForInput || session.State == SessionStateFailedSafe {
		if err := m.emitTrace(contract.RunID, contract.TraceID, ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(SessionStatePlanned, session.State, session.Reason, now)); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.observeMirror(&session); err != nil {
			return RuntimeResponse{}, err
		}
		if isTerminal(session.State) {
			if err := m.emitTerminalEvidence(&session); err != nil {
				return RuntimeResponse{}, err
			}
		}
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
		from := session.State
		if err := Transition(&session, SessionStateApproved, decision.Reason, now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.emitTrace(contract.RunID, contract.TraceID, ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateApproved, decision.Reason, now)); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.recordMemoryMutation(&session, "runtime.state", string(SessionStateApproved), decision.Reason); err != nil {
			return RuntimeResponse{}, err
		}
		if decision.AutoStart {
			from = session.State
			if err := Transition(&session, SessionStateRunning, "runtime session started", now); err != nil {
				return RuntimeResponse{}, err
			}
			if err := m.emitTrace(contract.RunID, contract.TraceID, ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateRunning, "runtime session started", now)); err != nil {
				return RuntimeResponse{}, err
			}
			if err := m.recordMemoryMutation(&session, "runtime.state", string(SessionStateRunning), "runtime session started"); err != nil {
				return RuntimeResponse{}, err
			}
		}
	case SessionStateWaitingForInput:
		from := session.State
		if err := Transition(&session, SessionStateWaitingForInput, decision.Reason, now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.emitTrace(contract.RunID, contract.TraceID, ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateWaitingForInput, decision.Reason, now)); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.recordMemoryMutation(&session, "runtime.state", string(SessionStateWaitingForInput), decision.Reason); err != nil {
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
		from := session.State
		if err := Transition(&session, SessionStateFailedSafe, decision.Reason, now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.emitTrace(contract.RunID, contract.TraceID, ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateFailedSafe, decision.Reason, now)); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.recordMemoryMutation(&session, "runtime.state", string(SessionStateFailedSafe), decision.Reason); err != nil {
			return RuntimeResponse{}, err
		}
	default:
		return RuntimeResponse{}, fmt.Errorf("unsupported admission initial state %s", decision.InitialState)
	}

	if err := m.observeMirror(&session); err != nil {
		return RuntimeResponse{}, err
	}
	if isTerminal(session.State) {
		if err := m.emitTerminalEvidence(&session); err != nil {
			return RuntimeResponse{}, err
		}
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
		from := session.State
		if err := Transition(&session, SessionStateRunning, "operator resume", now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateRunning, "operator resume", now)); err != nil {
			return RuntimeResponse{}, err
		}
	case SessionStateCheckpointed:
		from := session.State
		if err := Transition(&session, SessionStateRunning, "resume from checkpoint", now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateRunning, "resume from checkpoint", now)); err != nil {
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
		from := session.State
		if err := Transition(&session, SessionStateApproved, "operator approval acknowledged", now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateApproved, "operator approval acknowledged", now)); err != nil {
			return RuntimeResponse{}, err
		}
		from = session.State
		if err := Transition(&session, SessionStateRunning, "resume after approval", now); err != nil {
			return RuntimeResponse{}, err
		}
		if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(from, SessionStateRunning, "resume after approval", now)); err != nil {
			return RuntimeResponse{}, err
		}
	default:
		return RuntimeResponse{}, fmt.Errorf("cannot resume session in state %s", session.State)
	}

	if err := m.recordMemoryMutation(&session, "runtime.state", string(session.State), session.Reason); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.observeMirror(&session); err != nil {
		return RuntimeResponse{}, err
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
	if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(SessionStateRunning, SessionStatePaused, reason, now)); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.recordMemoryMutation(&session, "runtime.state", string(SessionStatePaused), reason); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.observeMirror(&session); err != nil {
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
	if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(SessionStateRunning, SessionStateCheckpointed, "checkpoint created", now)); err != nil {
		return RuntimeResponse{}, err
	}
	session.LatestCheckpointID = string(cp.Record.CheckpointID)
	if err := m.recordMemoryMutation(&session, "runtime.checkpoint", session.LatestCheckpointID, "checkpoint created"); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventArtifactLinked, map[string]any{"checkpoint_id": session.LatestCheckpointID}); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.observeMirror(&session); err != nil {
		return RuntimeResponse{}, err
	}

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
	if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventStateTransition, TransitionPayload(SessionStateCheckpointed, SessionStateRolledBack, "rollback to checkpoint "+checkpointID, now)); err != nil {
		return RuntimeResponse{}, err
	}
	session.LatestRollbackID = string(rb.Record.RollbackID)
	if err := m.recordMemoryMutation(&session, "runtime.rollback", session.LatestRollbackID, "rollback applied"); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.observeMirror(&session); err != nil {
		return RuntimeResponse{}, err
	}
	if err := m.emitTerminalEvidence(&session); err != nil {
		return RuntimeResponse{}, err
	}

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
		LatestMutationID:   session.LatestMutationID,
		LatestWitnessID:    session.LatestWitnessID,
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

func (m *Manager) Trace(context.Context, string) (map[string]any, error) {
	return nil, fmt.Errorf("trace inspect requires explicit subcommand input in W08")
}

func (m *Manager) Receipts(context.Context, string) (map[string]any, error) {
	return nil, fmt.Errorf("receipts inspect requires explicit subcommand input in W08")
}

func (m *Manager) Memory(_ context.Context, query string) (map[string]any, error) {
	if m.memory == nil {
		return nil, fmt.Errorf("memory plane is not wired")
	}
	return m.memory.Recall(query)
}

func (m *Manager) Mirror(_ context.Context, path string) (map[string]any, error) {
	if m.mirror == nil {
		return nil, fmt.Errorf("mirror plane is not wired")
	}
	return m.mirror.StatusFromFile(path)
}

func (m *Manager) TraceInspect(path string) (map[string]any, error) {
	if m.trace == nil {
		return nil, fmt.Errorf("trace plane is not wired")
	}
	return m.trace.InspectFromFile(path)
}

func (m *Manager) TraceReplay(runID string) (map[string]any, error) {
	if m.trace == nil {
		return nil, fmt.Errorf("trace plane is not wired")
	}
	result, err := m.trace.Replay.Replay(ids.RunID(runID))
	if err != nil {
		return nil, err
	}
	return map[string]any{"replay": result}, nil
}

func (m *Manager) ExportReceiptBundle(path string) (map[string]any, error) {
	if m.provenance == nil {
		return nil, fmt.Errorf("provenance plane is not wired")
	}
	return m.provenance.ExportBundle(path)
}

func (m *Manager) ClosureFromFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read closure input: %w", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("decode closure input: %w", err)
	}
	return payload, nil
}

func (m *Manager) recordMemoryMutation(session *RuntimeSession, key, value, reason string) error {
	if m.memory == nil {
		return nil
	}
	entries, err := m.memory.Ledger()
	if err != nil {
		return err
	}
	var parent ids.MemoryMutationID
	if len(entries) > 0 {
		parent = entries[len(entries)-1].MutationID
	}
	mutation, err := m.memory.Upsert(mem.MutationInput{
		Scope: mem.ScopeBinding{
			Scope:       "session",
			TenantID:    ids.TenantID(session.Contract.TenantID),
			WorkspaceID: ids.WorkspaceID(session.Contract.WorkspaceID),
			RunID:       ids.RunID(session.RunID),
			SessionID:   ids.SessionID(session.SessionID),
		},
		OperatorID:       ids.OperatorID(session.Contract.OperatorID),
		Kind:             "annotate",
		Branch:           branchForSession(*session, parent),
		ParentMutationID: parent,
		Key:              key,
		Value:            value,
		Reason:           reason,
		OccurredAt:       m.now().UTC(),
	})
	if err != nil {
		return err
	}
	session.LatestMutationID = string(mutation.MutationID)
	if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventMemoryMutation, map[string]any{
		"mutation_id": mutation.MutationID,
		"key":         mutation.Key,
		"lineage_ref": mutation.LineageRef,
	}); err != nil {
		return err
	}
	return nil
}

func branchForSession(session RuntimeSession, head ids.MemoryMutationID) cmemory.BranchReference {
	return cmemory.BranchReference{
		BranchID:       ids.BranchID("branch_" + session.SessionID),
		HeadMutationID: head,
	}
}

func (m *Manager) observeMirror(session *RuntimeSession) error {
	if m.mirror == nil || m.memory == nil {
		return nil
	}
	ledger, err := m.memory.Ledger()
	if err != nil {
		return err
	}
	branches, err := m.memory.Branches()
	if err != nil {
		return err
	}
	event, _, err := m.mirror.Observe(mirror.RuntimeSnapshot{
		RunID:           session.RunID,
		SessionID:       session.SessionID,
		RuntimeState:    string(session.State),
		PlanPath:        session.PlanPath,
		LastDecision:    session.LastDecision,
		MemoryMutations: len(ledger),
		BranchCount:     len(branches),
		PendingApproval: strings.TrimSpace(session.ApprovalPath) != "",
		Signals:         []string{session.Reason},
	}, "runtime_manager")
	if err != nil {
		return err
	}
	session.LatestWitnessID = event.WitnessID
	if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventMirrorWitness, map[string]any{
		"witness_id": event.WitnessID,
		"severity":   event.Severity,
		"basis":      event.Basis,
	}); err != nil {
		return err
	}
	return nil
}

func (m *Manager) emitTrace(runID ids.RunID, traceID ids.TraceID, spanID ids.SpanID, kind ctrace.EventType, payload map[string]any) error {
	if m.trace == nil || runID == "" || traceID == "" {
		return nil
	}
	_, err := m.trace.Emit(runID, traceID, spanID, ids.SpanID(""), kind, payload)
	return err
}

func (m *Manager) emitTerminalEvidence(session *RuntimeSession) error {
	if !isTerminal(session.State) || m.provenance == nil {
		return nil
	}
	traceRefs := []string{}
	if m.trace != nil {
		events, err := m.trace.Reader.Events(ids.RunID(session.RunID))
		if err != nil {
			return err
		}
		for _, event := range events {
			traceRefs = append(traceRefs, string(event.EventID))
		}
	}
	artifacts := []evidence.ArtifactRef{}
	if session.PlanPath != "" {
		if ref, err := m.provenance.LinkArtifact(ids.RunID(session.RunID), session.PlanPath, "application/json", "runtime_plan"); err == nil {
			artifacts = append(artifacts, ref)
		}
	}
	closure, closurePath, err := m.provenance.EmitClosure(ids.RunID(session.RunID), string(session.State), traceRefs, nil, nil)
	if err != nil {
		return err
	}
	if ref, err := m.provenance.LinkArtifact(ids.RunID(session.RunID), closurePath, "application/json", "closure_proof"); err == nil {
		artifacts = append(artifacts, ref)
	}
	receipt, err := m.provenance.EmitReceipt(ids.RunID(session.RunID), evidence.ReceiptKindClosure, string(session.State), "terminal run closure proof emitted", artifacts, closurePath)
	if err != nil {
		return err
	}
	if err := m.emitTrace(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), coretrace.EventClosureProof, map[string]any{
		"closure_status": closure.ClosureStatus,
		"closure_ref":    closurePath,
		"receipt_id":     receipt.ReceiptID,
	}); err != nil {
		return err
	}
	if m.trace != nil {
		event, err := m.trace.Emit(ids.RunID(session.RunID), ids.TraceID(session.Contract.TraceID), ids.SpanID(""), ids.SpanID(""), coretrace.EventReceiptLinked, map[string]any{"receipt_id": receipt.ReceiptID, "kind": receipt.Kind, "status": receipt.Status})
		if err != nil {
			return err
		}
		coretrace.LinkReceipt(&event, receipt.ReceiptID)
	}
	return nil
}

func isTerminal(state SessionState) bool {
	return state == SessionStateFailedSafe || state == SessionStateRolledBack || state == SessionStateCompleted
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
