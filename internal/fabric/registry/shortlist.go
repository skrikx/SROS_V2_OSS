package registry

import ctools "srosv2/contracts/tools"

func Shortlist(manifests []ctools.Manifest, query ctools.SearchQuery) []ctools.SearchMatch {
	q := query
	if q.Limit == 0 {
		q.Limit = 3
	}
	return Search(manifests, q)
}
