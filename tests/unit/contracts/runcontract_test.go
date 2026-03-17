package contracts_test

import (
	"encoding/json"
	"testing"
	"time"

	"srosv2/contracts/runcontract"
	"srosv2/internal/shared/ids"
)

func TestRunContractValidateValid(t *testing.T) {
	contract := runcontract.RunContract{
		ContractVersion:   "v2.0",
		RunID:             ids.RunID("run_001"),
		TraceID:           ids.TraceID("trace_001"),
		OperatorID:        ids.OperatorID("op_hassan"),
		TenantID:          ids.TenantID("local"),
		WorkspaceID:       ids.WorkspaceID("default"),
		IntentSummary:     "compile",
		NormalizedRequest: "go test ./...",
		RiskClass:         runcontract.RiskClassMedium,
		RouteClass:        runcontract.RouteClassLocalCLI,
		CreatedAt:         time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
	}

	errs := runcontract.Validate(contract)
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(errs))
	}
}

func TestRunContractValidateRejectsEnum(t *testing.T) {
	contract := runcontract.RunContract{
		ContractVersion:   "v2.0",
		RunID:             ids.RunID("run_001"),
		TraceID:           ids.TraceID("trace_001"),
		OperatorID:        ids.OperatorID("op_hassan"),
		TenantID:          ids.TenantID("local"),
		WorkspaceID:       ids.WorkspaceID("default"),
		IntentSummary:     "compile",
		NormalizedRequest: "go test ./...",
		RiskClass:         runcontract.RiskClass("invalid"),
		RouteClass:        runcontract.RouteClassLocalCLI,
		CreatedAt:         time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
	}

	errs := runcontract.Validate(contract)
	if len(errs) == 0 {
		t.Fatal("expected validation errors")
	}
}

func TestRunContractGoldenFixture(t *testing.T) {
	data := loadFixture(t, "run_contract.json")

	var contract runcontract.RunContract
	if err := json.Unmarshal(data, &contract); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	errs := runcontract.Validate(contract)
	if len(errs) != 0 {
		t.Fatalf("expected fixture to validate, got %d errors", len(errs))
	}

	raw, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal fixture roundtrip: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("expected non-empty roundtrip json")
	}
}
