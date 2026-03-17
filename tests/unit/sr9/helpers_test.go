package sr9_test

import (
	"time"

	"srosv2/contracts/runcontract"
)

func validContract() runcontract.RunContract {
	return runcontract.RunContract{
		ContractVersion:   "v2.0",
		RunID:             "run_sr9_001",
		TraceID:           "trace_sr9_001",
		OperatorID:        "op_local",
		TenantID:          "local",
		WorkspaceID:       "default",
		IntentSummary:     "sr9 test contract",
		NormalizedRequest: "sr9 test contract",
		RiskClass:         runcontract.RiskClassLow,
		RouteClass:        runcontract.RouteClassLocalCLI,
		RequestedReceipts: []runcontract.ReceiptRequest{{Mode: runcontract.ReceiptModeSummary, Kind: "compile"}},
		Metadata: map[string]string{
			"compile_request_id": "cmp_sr9_001",
			"domain_class":       "general_ops",
			"topology_class":     "local_single",
		},
		CreatedAt: time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC),
	}
}
