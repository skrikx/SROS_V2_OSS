package sr8_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"srosv2/contracts/evidence"
	"srosv2/contracts/runcontract"
	"srosv2/internal/core/sr8"
)

func TestBuildCompileReceipt(t *testing.T) {
	dir := t.TempDir()
	artifactPath := filepath.Join(dir, "run_contract.json")
	if err := os.WriteFile(artifactPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	ref, err := sr8.ArtifactRefFromFile(artifactPath)
	if err != nil {
		t.Fatalf("artifact ref from file: %v", err)
	}

	contract := runcontract.RunContract{RunID: "run_001", ContractVersion: "v2.0"}
	receipt, err := sr8.BuildCompileReceipt(
		"cmp_001",
		sr8.Classification{Domain: sr8.DomainGeneral, Risk: runcontract.RiskClassLow},
		sr8.TopologyDecision{Topology: sr8.TopologyLocalSingle, RouteClass: runcontract.RouteClassLocalCLI},
		contract,
		[]evidence.ArtifactRef{ref},
		time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("build compile receipt: %v", err)
	}
	if receipt.Status != "compiled" {
		t.Fatalf("unexpected receipt status: %s", receipt.Status)
	}
}
