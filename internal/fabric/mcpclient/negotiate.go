package mcpclient

import (
	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/registry"
)

func Negotiate(r *registry.Registry, req ctools.NegotiationRequest) ctools.NegotiationResult {
	matches := r.Shortlist(req.Query)
	if len(matches) == 0 {
		return ctools.NegotiationResult{Allowed: false, Reason: "no governed capability matched query"}
	}
	selected := matches[0]
	if !selected.Selectable && req.RequireRunnable {
		return ctools.NegotiationResult{Allowed: false, Alternates: matches, Reason: selected.Reason}
	}
	return ctools.NegotiationResult{
		Allowed:            selected.Selectable,
		Selected:           &selected,
		Alternates:         matches,
		Reason:             selected.Reason,
		PolicyFit:          selected.PolicyFit,
		EnvironmentFit:     selected.EnvironmentFit,
		SelectedCapability: selected.Manifest.Name,
	}
}
