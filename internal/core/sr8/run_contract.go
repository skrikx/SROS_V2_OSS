package sr8

import (
	"encoding/xml"
	"fmt"
	"time"

	"srosv2/contracts/runcontract"
	"srosv2/internal/shared/ids"
)

func AssembleRunContract(intent NormalizedIntent, class Classification, topo TopologyDecision, now time.Time) runcontract.RunContract {
	runID := ids.RunID("run_" + shortHash(intent.CompileRequestID+"|run"))
	traceID := ids.TraceID("trace_" + shortHash(intent.CompileRequestID+"|trace"))

	return runcontract.RunContract{
		ContractVersion:   "v2.0",
		RunID:             runID,
		TraceID:           traceID,
		OperatorID:        ids.OperatorID(intent.OperatorID),
		TenantID:          ids.TenantID(intent.TenantID),
		WorkspaceID:       ids.WorkspaceID(intent.WorkspaceID),
		IntentSummary:     intent.IntentSummary,
		NormalizedRequest: intent.NormalizedText,
		RiskClass:         class.Risk,
		RouteClass:        topo.RouteClass,
		RequestedReceipts: []runcontract.ReceiptRequest{{Mode: runcontract.ReceiptModeSummary, Kind: "compile"}},
		Metadata: map[string]string{
			"domain_class":       string(class.Domain),
			"topology_class":     string(topo.Topology),
			"artifact_posture":   class.ArtifactPosture,
			"compile_request_id": intent.CompileRequestID,
		},
		CreatedAt: now.UTC(),
	}
}

func RenderSRXML(contract runcontract.RunContract) (string, error) {
	type runContractXML struct {
		XMLName           xml.Name `xml:"run_contract"`
		ContractVersion   string   `xml:"contract_version,attr"`
		RunID             string   `xml:"run_id"`
		TraceID           string   `xml:"trace_id"`
		OperatorID        string   `xml:"operator_id"`
		TenantID          string   `xml:"tenant_id"`
		WorkspaceID       string   `xml:"workspace_id"`
		IntentSummary     string   `xml:"intent_summary"`
		NormalizedRequest string   `xml:"normalized_request"`
		RiskClass         string   `xml:"risk_class"`
		RouteClass        string   `xml:"route_class"`
		CreatedAt         string   `xml:"created_at"`
	}
	type srxmlDoc struct {
		XMLName     xml.Name       `xml:"srxml"`
		Version     string         `xml:"version,attr"`
		RunContract runContractXML `xml:"run_contract"`
	}

	doc := srxmlDoc{
		Version: "sros.v2",
		RunContract: runContractXML{
			ContractVersion:   contract.ContractVersion,
			RunID:             string(contract.RunID),
			TraceID:           string(contract.TraceID),
			OperatorID:        string(contract.OperatorID),
			TenantID:          string(contract.TenantID),
			WorkspaceID:       string(contract.WorkspaceID),
			IntentSummary:     contract.IntentSummary,
			NormalizedRequest: contract.NormalizedRequest,
			RiskClass:         string(contract.RiskClass),
			RouteClass:        string(contract.RouteClass),
			CreatedAt:         contract.CreatedAt.UTC().Format(time.RFC3339),
		},
	}

	bytes, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal srxml: %w", err)
	}
	return xml.Header + string(bytes), nil
}
