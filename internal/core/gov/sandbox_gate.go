package gov

import (
	"fmt"
	"strings"

	"srosv2/contracts/policy"
)

func resolveSandbox(bundle policy.Bundle, capability string) (policy.SandboxProfile, error) {
	rule, ok := matchCapability(bundle, capability)
	name := bundle.DefaultSandboxProfile
	if ok && rule.SandboxProfile != "" {
		name = rule.SandboxProfile
	}
	if strings.TrimSpace(name) == "" {
		return policy.SandboxProfile{}, fmt.Errorf("capability %s has no sandbox profile", capability)
	}
	profile, ok := bundle.Sandboxes[name]
	if !ok {
		return policy.SandboxProfile{}, fmt.Errorf("sandbox profile %s not declared", name)
	}
	return profile, nil
}

func sandboxAllows(profile policy.SandboxProfile, capability string) bool {
	switch {
	case strings.HasPrefix(capability, "shell."):
		return profile.AllowShell
	case strings.HasPrefix(capability, "patch."):
		return profile.AllowPatch
	case strings.HasPrefix(capability, "connector."), strings.HasPrefix(capability, "mcp."):
		return profile.AllowExternalNet
	default:
		return true
	}
}
