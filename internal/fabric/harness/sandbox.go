package harness

import "strings"

func Compatible(profile Profile, capabilityClass string) bool {
	switch {
	case strings.Contains(capabilityClass, ".patch"):
		return profile.AllowPatch
	case strings.Contains(capabilityClass, ".shell"):
		return profile.AllowShell
	case strings.HasPrefix(capabilityClass, "connector"), strings.Contains(capabilityClass, ".web"), strings.HasPrefix(capabilityClass, "mcp"):
		return profile.AllowNetwork || profile.Name == "network_restricted"
	default:
		return profile.AllowRead || profile.AllowWrite || profile.AllowPatch || profile.AllowShell
	}
}
