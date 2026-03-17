package sr8

import (
	"srosv2/internal/ace/intake"
)

func SelectTopology(intent NormalizedIntent, class Classification) TopologyDecision {
	decision := intake.SelectTopology(intake.Classification{
		Domain: intake.DomainClass(class.Domain),
		Risk:   class.Risk,
	})

	return TopologyDecision{
		Topology:   TopologyClass(decision.Topology),
		RouteClass: decision.RouteClass,
		Reason:     decision.Reason,
	}
}
