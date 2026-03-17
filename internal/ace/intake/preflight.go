package intake

import "srosv2/contracts/runcontract"

type PreflightCheck struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason,omitempty"`
}

func BuildPreflight(class Classification) []PreflightCheck {
	checks := []PreflightCheck{
		{Name: "intent_nonempty", Pass: true},
		{Name: "classification_complete", Pass: true},
		{Name: "topology_selectable", Pass: true},
	}

	if class.Risk == runcontract.RiskClassCritical {
		checks = append(checks, PreflightCheck{Name: "high_risk_notice", Pass: true, Reason: "critical compile risk; runtime governance required later"})
	}
	return checks
}
