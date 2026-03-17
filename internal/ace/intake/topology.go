package intake

import "srosv2/contracts/runcontract"

type TopologyClass string

const (
	TopologyLocalSingle     TopologyClass = "local_single"
	TopologyLocalFilesystem TopologyClass = "local_filesystem"
	TopologyLocalResearch   TopologyClass = "local_research"
)

type TopologyDecision struct {
	Topology   TopologyClass          `json:"topology"`
	RouteClass runcontract.RouteClass `json:"route_class"`
	Reason     string                 `json:"reason"`
}

func SelectTopology(class Classification) TopologyDecision {
	switch class.Domain {
	case DomainFileTask:
		return TopologyDecision{
			Topology:   TopologyLocalFilesystem,
			RouteClass: runcontract.RouteClassLocalRun,
			Reason:     "file task classified at compile edge",
		}
	case DomainResearch:
		return TopologyDecision{
			Topology:   TopologyLocalResearch,
			RouteClass: runcontract.RouteClassLocalCLI,
			Reason:     "research task remains local inspection oriented",
		}
	default:
		return TopologyDecision{
			Topology:   TopologyLocalSingle,
			RouteClass: runcontract.RouteClassLocalCLI,
			Reason:     "general local compile route",
		}
	}
}
