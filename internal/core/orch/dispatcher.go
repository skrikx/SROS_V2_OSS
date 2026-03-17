package orch

import (
	"context"
	"fmt"
	"time"
)

type Decision struct {
	Verdict              string
	Reason               string
	SandboxProfile       string
	ApprovalCheckpointID string
}

type DecisionFunc func(context.Context, WorkUnit) (Decision, error)

type ExecutionResult struct {
	LastUnit     string           `json:"last_unit,omitempty"`
	Decision     Decision         `json:"decision"`
	Route        *CheckpointRoute `json:"route,omitempty"`
	Completed    bool             `json:"completed"`
	ExecutedUnit []string         `json:"executed_units"`
}

type Dispatcher struct {
	bus    *Bus
	router *CheckpointRouter
	now    func() time.Time
}

func NewDispatcher(bus *Bus, router *CheckpointRouter, now func() time.Time) *Dispatcher {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Dispatcher{bus: bus, router: router, now: now}
}

func (d *Dispatcher) Execute(ctx context.Context, plan Plan, queue *Queue, decide DecisionFunc) (ExecutionResult, error) {
	if plan.Concurrency.MaxParallel < 1 {
		return ExecutionResult{}, fmt.Errorf("invalid concurrency rule")
	}
	result := ExecutionResult{Completed: true, ExecutedUnit: []string{}}

	for queue.Len() > 0 {
		select {
		case <-ctx.Done():
			return ExecutionResult{}, ctx.Err()
		default:
		}

		unit, err := queue.Dequeue()
		if err != nil {
			return ExecutionResult{}, err
		}
		result.LastUnit = unit.ID
		result.ExecutedUnit = append(result.ExecutedUnit, unit.ID)
		if d.bus != nil {
			d.bus.Publish(Event{Type: "work_unit.started", SessionID: plan.SessionID, WorkUnit: unit.ID, Message: unit.Name, At: d.now()})
		}

		if !unit.Privileged {
			result.Decision = Decision{Verdict: "allow", Reason: "non-privileged work unit"}
			continue
		}

		decision, err := decide(ctx, unit)
		if err != nil {
			return ExecutionResult{}, err
		}
		result.Decision = decision

		if d.bus != nil {
			d.bus.Publish(Event{Type: "work_unit.governed", SessionID: plan.SessionID, WorkUnit: unit.ID, Message: decision.Verdict + ": " + decision.Reason, At: d.now()})
		}

		switch decision.Verdict {
		case "allow":
		case "ask":
			result.Completed = false
			if d.router == nil {
				return ExecutionResult{}, fmt.Errorf("checkpoint router is not configured")
			}
			route, err := d.router.RouteAsk(CheckpointRoute{
				SessionID:      plan.SessionID,
				WorkUnitID:     unit.ID,
				Route:          "local_cli_checkpoint",
				RequestedAt:    d.now(),
				Reason:         decision.Reason,
				Capability:     unit.Capability,
				SandboxProfile: decision.SandboxProfile,
			})
			if err != nil {
				return ExecutionResult{}, err
			}
			result.Route = &route
			return result, nil
		case "deny":
			result.Completed = false
			return result, nil
		default:
			return ExecutionResult{}, fmt.Errorf("unsupported verdict %q", decision.Verdict)
		}
	}

	return result, nil
}
