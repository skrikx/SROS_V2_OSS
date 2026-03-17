package runtime_integration_test

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

func writeRuntimeContract(t *testing.T, dir string, metadata map[string]string) string {
	t.Helper()
	contract := runcontract.RunContract{
		ContractVersion:   "v2.0",
		RunID:             "run_001",
		TraceID:           "trace_001",
		OperatorID:        "op_local",
		TenantID:          "local",
		WorkspaceID:       "default",
		IntentSummary:     "governed runtime session",
		NormalizedRequest: "governed runtime session",
		RiskClass:         runcontract.RiskClassHigh,
		RouteClass:        runcontract.RouteClassLocalCLI,
		Metadata:          metadata,
		CreatedAt:         fixedNow,
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

var fixedNow = time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)

type stubGate struct {
	decision runtime.AdmissionDecision
}

func (s stubGate) Admit(context.Context, runtime.AdmissionRequest) (runtime.AdmissionDecision, error) {
	return s.decision, nil
}
