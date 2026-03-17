package orch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/contracts/runcontract"
)

type Options struct {
	ArtifactRoot string
	Now          func() time.Time
	EventHook    func(string, map[string]any)
}

type Orchestrator struct {
	artifactRoot string
	scheduler    *Scheduler
	bus          *Bus
	router       *CheckpointRouter
	now          func() time.Time
	eventHook    func(string, map[string]any)
}

func New(opts Options) (*Orchestrator, error) {
	root := opts.ArtifactRoot
	if root == "" {
		root = filepath.Join("artifacts", "runtime", "orch")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create orchestrator artifact root: %w", err)
	}
	router, err := NewCheckpointRouter(filepath.Join(filepath.Dir(root), "approvals"))
	if err != nil {
		return nil, err
	}
	now := opts.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Orchestrator{
		artifactRoot: root,
		scheduler:    NewScheduler(),
		bus:          NewBus(),
		router:       router,
		now:          now,
		eventHook:    opts.EventHook,
	}, nil
}

func (o *Orchestrator) Hydrate(sessionID string, contract runcontract.RunContract, topologyHint string) (Plan, error) {
	plan, err := o.scheduler.Build(sessionID, contract, topologyHint)
	if err != nil {
		return Plan{}, err
	}
	if err := o.writeJSON(sessionID+"_plan.json", plan); err != nil {
		return Plan{}, err
	}
	if o.bus != nil {
		o.bus.Publish(Event{Type: "plan.hydrated", SessionID: sessionID, Message: "orchestration plan hydrated", At: o.now()})
	}
	if o.eventHook != nil {
		o.eventHook("plan.hydrated", map[string]any{"session_id": sessionID, "run_id": string(contract.RunID)})
	}
	return plan, nil
}

func (o *Orchestrator) Execute(ctx context.Context, plan Plan, decide DecisionFunc) (ExecutionResult, error) {
	queue := NewQueue()
	for _, unit := range plan.WorkUnits {
		queue.Enqueue(unit)
	}
	dispatcher := NewDispatcher(o.bus, o.router, o.now)
	result, err := dispatcher.Execute(ctx, plan, queue, decide)
	if err != nil {
		return ExecutionResult{}, err
	}
	if err := o.writeJSON(plan.SessionID+"_execution.json", result); err != nil {
		return ExecutionResult{}, err
	}
	if err := o.writeJSON(plan.SessionID+"_events.json", o.bus.Events()); err != nil {
		return ExecutionResult{}, err
	}
	if o.eventHook != nil {
		o.eventHook("plan.executed", map[string]any{"session_id": plan.SessionID, "run_id": plan.RunID, "completed": result.Completed})
	}
	return result, nil
}

func (o *Orchestrator) SetEventHook(hook func(string, map[string]any)) {
	o.eventHook = hook
}

func (o *Orchestrator) writeJSON(name string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal orchestrator artifact: %w", err)
	}
	path := filepath.Join(o.artifactRoot, name)
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write orchestrator artifact: %w", err)
	}
	return nil
}
