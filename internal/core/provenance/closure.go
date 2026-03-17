package provenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/internal/shared/ids"
)

type ClosureProof struct {
	RunID         ids.RunID `json:"run_id"`
	TerminalState string    `json:"terminal_state"`
	TraceRefs     []string  `json:"trace_refs"`
	ReceiptRefs   []string  `json:"receipt_refs"`
	ArtifactRefs  []string  `json:"artifact_refs"`
	ClosureStatus string    `json:"closure_status"`
	GeneratedAt   time.Time `json:"generated_at"`
}

func (s *Service) EmitClosure(runID ids.RunID, terminalState string, traceRefs, receiptRefs, artifactRefs []string) (ClosureProof, string, error) {
	proof := ClosureProof{
		RunID:         runID,
		TerminalState: terminalState,
		TraceRefs:     traceRefs,
		ReceiptRefs:   receiptRefs,
		ArtifactRefs:  artifactRefs,
		ClosureStatus: "sealed",
		GeneratedAt:   s.now().UTC(),
	}
	ref := filepath.Join(s.root, "closures", string(runID)+"_closure.json")
	data, err := json.MarshalIndent(proof, "", "  ")
	if err != nil {
		return ClosureProof{}, "", fmt.Errorf("marshal closure proof: %w", err)
	}
	if err := os.WriteFile(ref, append(data, '\n'), 0o644); err != nil {
		return ClosureProof{}, "", fmt.Errorf("write closure proof: %w", err)
	}
	return proof, ref, nil
}
