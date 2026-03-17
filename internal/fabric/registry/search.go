package registry

import (
	"sort"
	"strings"

	ctools "srosv2/contracts/tools"
)

func Search(manifests []ctools.Manifest, query ctools.SearchQuery) []ctools.SearchMatch {
	results := []ctools.SearchMatch{}
	for _, manifest := range manifests {
		match, ok := scoreManifest(manifest, query)
		if ok {
			results = append(results, match)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Rank == results[j].Rank {
			return results[i].Manifest.Name < results[j].Manifest.Name
		}
		return results[i].Rank > results[j].Rank
	})
	if query.Limit > 0 && len(results) > query.Limit {
		return results[:query.Limit]
	}
	return results
}

func scoreManifest(manifest ctools.Manifest, query ctools.SearchQuery) (ctools.SearchMatch, bool) {
	if query.Class != "" && !strings.HasPrefix(manifest.Class, query.Class) {
		return ctools.SearchMatch{}, false
	}
	if query.Domain != "" && manifest.Domain != query.Domain {
		return ctools.SearchMatch{}, false
	}
	if query.PolicyClass != "" && manifest.PolicyClass != query.PolicyClass {
		return ctools.SearchMatch{}, false
	}
	if query.AuthType != "" && manifest.AuthType != query.AuthType {
		return ctools.SearchMatch{}, false
	}
	if query.TrustBoundary != "" && manifest.TrustBoundary != query.TrustBoundary {
		return ctools.SearchMatch{}, false
	}
	if !query.IncludeUnsafe && manifest.Unsafe {
		return ctools.SearchMatch{}, false
	}
	if len(query.Status) > 0 {
		allowed := false
		for _, state := range query.Status {
			if manifest.Status == state {
				allowed = true
				break
			}
		}
		if !allowed {
			return ctools.SearchMatch{}, false
		}
	}
	if query.MCPMode == "ingress" && !manifest.MCPIngressCapable {
		return ctools.SearchMatch{}, false
	}
	if query.MCPMode == "egress" && !manifest.MCPEgressCapable {
		return ctools.SearchMatch{}, false
	}

	rank := 10
	if manifest.Status == ctools.StateActive {
		rank += 5
	}
	if manifest.Status == ctools.StateAdmitted || manifest.Status == ctools.StateExperimental {
		rank += 3
	}
	selectable := manifest.Status == ctools.StateActive || manifest.Status == ctools.StateExperimental || manifest.Status == ctools.StateAdmitted
	reason := "policy and environment compatible"
	if manifest.Status == ctools.StateQuarantined {
		selectable = false
		reason = "quarantined capabilities remain visible but not selectable"
	}
	if manifest.Status == ctools.StateDeprecated {
		selectable = false
		reason = "deprecated capability remains historically visible"
	}
	if manifest.Status == ctools.StateDisabled || manifest.Status == ctools.StateDraft || manifest.Status == ctools.StateValidated {
		selectable = false
		reason = "capability is not runnable in current lifecycle state"
	}

	return ctools.SearchMatch{
		Rank:                 rank,
		Manifest:             manifest.Summary(),
		PolicyFit:            "policy_bound",
		EnvironmentFit:       "local_safe",
		AuthEnvelopeRequired: manifest.AuthEnvelopeRequired,
		SelectedSandbox:      firstSandbox(manifest.SandboxProfiles),
		Unsafe:               manifest.Unsafe,
		Selectable:           selectable,
		Reason:               reason,
	}, true
}

func firstSandbox(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
