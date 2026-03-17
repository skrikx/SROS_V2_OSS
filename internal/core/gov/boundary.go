package gov

import (
	"strings"

	"srosv2/contracts/policy"
)

func ResolveBoundary(capability string) policy.TrustBoundary {
	switch {
	case strings.HasPrefix(capability, "patch."):
		return policy.TrustBoundaryLocalFS
	case strings.HasPrefix(capability, "connector."), strings.HasPrefix(capability, "mcp."):
		return policy.TrustBoundaryExternalNet
	default:
		return policy.TrustBoundaryLocalProcess
	}
}
