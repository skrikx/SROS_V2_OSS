package gov

import (
	"strings"

	"srosv2/contracts/policy"
	"srosv2/contracts/runcontract"
)

func isPrivilegedCapability(capability string) bool {
	return strings.HasPrefix(capability, "shell.") ||
		strings.HasPrefix(capability, "patch.") ||
		strings.HasPrefix(capability, "tool.") ||
		strings.HasPrefix(capability, "connector.") ||
		strings.HasPrefix(capability, "mcp.")
}

func permissionVerdict(bundle policy.Bundle, capability string, risk runcontract.RiskClass) policy.Verdict {
	rule, ok := matchCapability(bundle, capability)
	if ok {
		if rule.RequireApproval {
			return policy.VerdictAsk
		}
		return rule.Verdict
	}
	if strings.HasPrefix(capability, "shell.") || strings.HasPrefix(capability, "patch.") {
		return policy.VerdictDeny
	}
	if risk == runcontract.RiskClassCritical || risk == runcontract.RiskClassHigh {
		return policy.VerdictAsk
	}
	if bundle.DefaultVerdict != "" {
		return bundle.DefaultVerdict
	}
	if isPrivilegedCapability(capability) {
		return policy.VerdictAsk
	}
	return policy.VerdictAllow
}
