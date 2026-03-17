package sr8_test

import (
	"testing"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/sr8"
)

func TestSelectTopology(t *testing.T) {
	decision := sr8.SelectTopology(sr8.NormalizedIntent{NormalizedText: "patch file now"}, sr8.Classification{
		Domain: sr8.DomainFileTask,
		Risk:   runcontract.RiskClassMedium,
	})

	if decision.Topology != sr8.TopologyLocalFilesystem {
		t.Fatalf("unexpected topology: %s", decision.Topology)
	}
	if decision.RouteClass != runcontract.RouteClassLocalRun {
		t.Fatalf("unexpected route class: %s", decision.RouteClass)
	}
}
