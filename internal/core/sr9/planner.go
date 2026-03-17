package sr9

import (
	"srosv2/contracts/policy"
	"srosv2/contracts/runcontract"
)

type Planner interface {
	Decide(runcontract.RunContract, Binding) (policy.Verdict, string)
}

type defaultPlanner struct{}

func (defaultPlanner) Decide(contract runcontract.RunContract, binding Binding) (policy.Verdict, string) {
	if contract.Metadata["approval_mode"] == "deny" {
		return policy.VerdictDeny, "runtime preflight requested explicit deny mode"
	}
	if contract.Metadata["approval_mode"] == "ask" {
		return policy.VerdictAsk, "runtime preflight requires local operator approval checkpoint"
	}
	if contract.RiskClass == runcontract.RiskClassCritical {
		return policy.VerdictAsk, "critical risk requires local operator checkpoint"
	}
	return policy.VerdictAllow, "runtime preflight allow for " + binding.RuntimeShell
}
