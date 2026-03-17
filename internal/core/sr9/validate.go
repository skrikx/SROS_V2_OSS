package sr9

import (
	"fmt"

	"srosv2/contracts/runcontract"
)

func ValidateContract(contract runcontract.RunContract) error {
	errs := runcontract.Validate(contract)
	if len(errs) > 0 {
		return fmt.Errorf("run contract validation failed: %v", errs[0])
	}
	if contract.Metadata == nil {
		return fmt.Errorf("run contract metadata is required for runtime admission")
	}
	if contract.Metadata["compile_request_id"] == "" {
		return fmt.Errorf("run contract missing compile_request_id")
	}
	if contract.Metadata["domain_class"] == "" {
		return fmt.Errorf("run contract missing domain_class")
	}
	if contract.Metadata["topology_class"] == "" {
		return fmt.Errorf("run contract missing topology_class")
	}
	return nil
}
