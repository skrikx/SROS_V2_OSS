package sr8

import (
	"fmt"
	"strings"

	"srosv2/contracts/runcontract"
)

func ValidateCompileInput(intent NormalizedIntent, class Classification, topo TopologyDecision) error {
	if strings.TrimSpace(intent.NormalizedText) == "" {
		return fmt.Errorf("normalized intent is empty")
	}
	if strings.TrimSpace(intent.CompileRequestID) == "" {
		return fmt.Errorf("compile_request_id is required")
	}
	if class.Domain == "" {
		return fmt.Errorf("classification domain is required")
	}
	if class.Risk == "" {
		return fmt.Errorf("classification risk is required")
	}
	if topo.Topology == "" || topo.RouteClass == "" {
		return fmt.Errorf("topology decision is incomplete")
	}
	return nil
}

func ValidateRunContract(contract runcontract.RunContract) error {
	errs := runcontract.Validate(contract)
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("run contract validation failed: %v", errs[0])
}
