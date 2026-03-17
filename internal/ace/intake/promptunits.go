package intake

import "srosv2/contracts/runcontract"

type PromptUnitHint struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func SuggestPromptUnits(domain DomainClass, risk runcontract.RiskClass) []PromptUnitHint {
	hints := []PromptUnitHint{}
	switch domain {
	case DomainFileTask:
		hints = append(hints, PromptUnitHint{Name: "file_delta_plan", Reason: "structured patch task"})
	case DomainResearch:
		hints = append(hints, PromptUnitHint{Name: "research_scope", Reason: "bounded research framing"})
	default:
		hints = append(hints, PromptUnitHint{Name: "operator_intent", Reason: "default compile framing"})
	}

	if risk == runcontract.RiskClassHigh || risk == runcontract.RiskClassCritical {
		hints = append(hints, PromptUnitHint{Name: "risk_notice", Reason: "high-risk compile classification"})
	}

	return hints
}
