package gov

import "srosv2/contracts/policy"

func applyBreakGlass(bundle policy.Bundle, rule *policy.CapabilityPolicy, requested bool) (policy.Verdict, string) {
	if !requested {
		return "", ""
	}
	if !bundle.BreakGlassAllowed {
		return policy.VerdictDeny, "break-glass requested but bundle forbids it"
	}
	if rule != nil && !rule.AllowBreakGlass {
		return policy.VerdictDeny, "break-glass requested but capability forbids it"
	}
	return policy.VerdictAsk, "break-glass requested and requires explicit local approval"
}
