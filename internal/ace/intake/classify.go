package intake

import (
	"strings"

	"srosv2/contracts/runcontract"
)

type DomainClass string

const (
	DomainGeneral  DomainClass = "general_ops"
	DomainFileTask DomainClass = "file_patch"
	DomainResearch DomainClass = "research"
)

type Classification struct {
	Domain  DomainClass           `json:"domain"`
	Risk    runcontract.RiskClass `json:"risk"`
	Signals []string              `json:"signals,omitempty"`
}

func ClassifyIntent(text string) Classification {
	t := strings.ToLower(text)
	signals := make([]string, 0)

	domain := DomainGeneral
	switch {
	case containsAny(t, "patch", "file", "edit", "diff", "refactor"):
		domain = DomainFileTask
		signals = append(signals, "file_operations")
	case containsAny(t, "research", "investigate", "analyze", "report"):
		domain = DomainResearch
		signals = append(signals, "research_language")
	}

	risk := runcontract.RiskClassLow
	switch {
	case containsAny(t, "rm -rf", "delete all", "drop table", "format disk"):
		risk = runcontract.RiskClassCritical
		signals = append(signals, "destructive_keywords")
	case containsAny(t, "delete", "rollback", "production", "privileged", "sudo"):
		risk = runcontract.RiskClassHigh
		signals = append(signals, "high_risk_keywords")
	case containsAny(t, "network", "remote", "publish", "download"):
		risk = runcontract.RiskClassMedium
		signals = append(signals, "external_surface_keywords")
	}

	return Classification{Domain: domain, Risk: risk, Signals: signals}
}

func containsAny(v string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(v, n) {
			return true
		}
	}
	return false
}
