package sr8_test

import (
	"strings"
	"testing"
	"time"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/sr8"
)

func TestAssembleRunContract(t *testing.T) {
	norm := sr8.NormalizedIntent{
		NormalizedText:   "patch file",
		IntentSummary:    "patch file",
		OperatorID:       "op_local",
		TenantID:         "local",
		WorkspaceID:      "default",
		CompileRequestID: "cmp_abc123",
	}
	class := sr8.Classification{Domain: sr8.DomainFileTask, Risk: runcontract.RiskClassMedium, ArtifactPosture: "file_delta"}
	topo := sr8.TopologyDecision{Topology: sr8.TopologyLocalFilesystem, RouteClass: runcontract.RouteClassLocalRun, Reason: "file task"}

	contract := sr8.AssembleRunContract(norm, class, topo, time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC))
	if err := sr8.ValidateRunContract(contract); err != nil {
		t.Fatalf("run contract should validate: %v", err)
	}

	srxml, err := sr8.RenderSRXML(contract)
	if err != nil {
		t.Fatalf("render srxml: %v", err)
	}
	if !strings.Contains(srxml, "<run_contract") {
		t.Fatalf("expected run_contract element in srxml: %s", srxml)
	}
}
