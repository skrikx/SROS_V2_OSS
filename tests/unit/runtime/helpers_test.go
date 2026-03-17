package runtime_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/runtime"
)

var fixedNow = time.Date(2026, 3, 17, 15, 4, 5, 0, time.UTC)

func writeRunContractFile(t *testing.T, dir string, mutate func(*runcontract.RunContract)) string {
	t.Helper()

	contract := runcontract.RunContract{
		ContractVersion:   "v2.0",
		RunID:             "run_test_001",
		TraceID:           "trace_test_001",
		OperatorID:        "op_local",
		TenantID:          "local",
		WorkspaceID:       "default",
		IntentSummary:     "runtime test contract",
		NormalizedRequest: "runtime test contract",
		RiskClass:         runcontract.RiskClassLow,
		RouteClass:        runcontract.RouteClassLocalCLI,
		RequestedReceipts: []runcontract.ReceiptRequest{{Mode: runcontract.ReceiptModeSummary, Kind: "compile"}},
		Metadata: map[string]string{
			"compile_request_id": "cmp_test_001",
			"domain_class":       "general_ops",
			"topology_class":     "local_single",
		},
		CreatedAt: fixedNow,
	}
	if mutate != nil {
		mutate(&contract)
	}

	path := filepath.Join(dir, "run_contract.json")
	data, err := json.MarshalIndent(contract, "", "  ")
	if err != nil {
		t.Fatalf("marshal contract: %v", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}
	return path
}

type stubGate struct {
	decision runtime.AdmissionDecision
	err      error
}

func (s stubGate) Admit(_ context.Context, _ runtime.AdmissionRequest) (runtime.AdmissionDecision, error) {
	if s.err != nil {
		return runtime.AdmissionDecision{}, s.err
	}
	return s.decision, nil
}
