package sr9_test

import (
	"testing"

	"srosv2/internal/core/sr9"
)

func TestValidateContractPassesForCanonicalInput(t *testing.T) {
	if err := sr9.ValidateContract(validContract()); err != nil {
		t.Fatalf("validate contract: %v", err)
	}
}

func TestValidateContractRequiresMetadata(t *testing.T) {
	contract := validContract()
	contract.Metadata = nil
	if err := sr9.ValidateContract(contract); err == nil {
		t.Fatal("expected validation error for missing metadata")
	}
}
