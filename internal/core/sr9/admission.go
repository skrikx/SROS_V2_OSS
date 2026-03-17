package sr9

import (
	"fmt"

	"srosv2/contracts/policy"
	"srosv2/contracts/runcontract"
)

type Admission struct {
	Verdict policy.Verdict `json:"verdict"`
	Reason  string         `json:"reason"`
	Binding Binding        `json:"binding"`
}

func BuildAdmission(contract runcontract.RunContract, planner Planner) (Admission, error) {
	if err := ValidateContract(contract); err != nil {
		return Admission{}, err
	}

	binding := BindTopology(contract)
	if planner == nil {
		planner = defaultPlanner{}
	}

	verdict, reason := planner.Decide(contract, binding)
	switch verdict {
	case policy.VerdictAllow, policy.VerdictAsk, policy.VerdictDeny:
		return Admission{Verdict: verdict, Reason: reason, Binding: binding}, nil
	default:
		return Admission{}, fmt.Errorf("unsupported admission verdict %q", verdict)
	}
}
