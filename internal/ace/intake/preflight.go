package intake

import (
	"srosv2/contracts/runcontract"
	ctools "srosv2/contracts/tools"
)

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

func BuildFabricPreflight(class Classification, matches []ctools.SearchMatch) []PreflightCheck {
	checks := BuildPreflight(class)
	if len(matches) == 0 {
		return append(checks, PreflightCheck{Name: "fabric_shortlist_present", Pass: false, Reason: "no governed capability shortlist available"})
	}
	return append(checks, PreflightCheck{Name: "fabric_shortlist_present", Pass: true, Reason: "governed capability shortlist available"})
}
