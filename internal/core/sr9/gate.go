package sr9

import (
	"context"

	"srosv2/contracts/policy"
	"srosv2/internal/core/runtime"
)

type Gate struct {
	planner Planner
}

type Options struct {
	Planner Planner
}

func NewGate(opts Options) *Gate {
	planner := opts.Planner
	if planner == nil {
		planner = defaultPlanner{}
	}
	return &Gate{planner: planner}
}

func (g *Gate) Admit(_ context.Context, req runtime.AdmissionRequest) (runtime.AdmissionDecision, error) {
	admission, err := BuildAdmission(req.Contract, g.planner)
	if err != nil {
		return runtime.AdmissionDecision{}, err
	}

	switch admission.Verdict {
	case policy.VerdictAllow:
		return runtime.AdmissionDecision{
			InitialState:    runtime.SessionStateApproved,
			Reason:          admission.Reason,
			AutoStart:       true,
			TopologyBinding: admission.Binding.RuntimeShell,
		}, nil
	case policy.VerdictAsk:
		return runtime.AdmissionDecision{
			InitialState:        runtime.SessionStateWaitingForInput,
			Reason:              admission.Reason,
			AutoStart:           false,
			RequireOperatorAck:  true,
			WaitingApprovalHint: "create local approval artifact with {\"approved\":true}",
			TopologyBinding:     admission.Binding.RuntimeShell,
		}, nil
	default:
		return runtime.AdmissionDecision{
			InitialState:    runtime.SessionStateFailedSafe,
			Reason:          admission.Reason,
			AutoStart:       false,
			TopologyBinding: admission.Binding.RuntimeShell,
		}, nil
	}
}
