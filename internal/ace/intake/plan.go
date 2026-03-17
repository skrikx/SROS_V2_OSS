package intake

type Plan struct {
	Classification  Classification   `json:"classification"`
	Topology        TopologyDecision `json:"topology"`
	Shortlist       Shortlist        `json:"shortlist"`
	Preflight       []PreflightCheck `json:"preflight"`
	ArtifactPosture string           `json:"artifact_posture"`
}

func BuildPlan(intent string) Plan {
	class := ClassifyIntent(intent)
	topo := SelectTopology(class)
	shortlist := BuildShortlist(class)
	preflight := BuildPreflight(class)

	posture := "run_contract_only"
	switch class.Domain {
	case DomainFileTask:
		posture = "file_delta"
	case DomainResearch:
		posture = "research_brief"
	}

	return Plan{
		Classification:  class,
		Topology:        topo,
		Shortlist:       shortlist,
		Preflight:       preflight,
		ArtifactPosture: posture,
	}
}
