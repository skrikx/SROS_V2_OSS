package sr9_test

import (
	"testing"

	"srosv2/internal/core/sr9"
)

func TestBindTopologyUsesContractMetadata(t *testing.T) {
	contract := validContract()
	contract.Metadata["topology_class"] = "local_filesystem"

	binding := sr9.BindTopology(contract)
	if binding.TopologyClass != "local_filesystem" {
		t.Fatalf("unexpected topology class %s", binding.TopologyClass)
	}
	if binding.RuntimeShell != "local_fs_session" {
		t.Fatalf("unexpected runtime shell %s", binding.RuntimeShell)
	}
}
