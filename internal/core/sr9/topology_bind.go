package sr9

import "srosv2/contracts/runcontract"

type Binding struct {
	TopologyClass string `json:"topology_class"`
	RuntimeShell  string `json:"runtime_shell"`
}

func BindTopology(contract runcontract.RunContract) Binding {
	topology := contract.Metadata["topology_class"]
	if topology == "" {
		topology = string(contract.RouteClass)
	}

	runtimeShell := "local_cli_session"
	switch topology {
	case "local_filesystem":
		runtimeShell = "local_fs_session"
	case "local_research":
		runtimeShell = "local_research_session"
	}

	return Binding{TopologyClass: topology, RuntimeShell: runtimeShell}
}
