package evidence

import (
	"time"

	"srosv2/internal/shared/ids"
)

type ReceiptKind string

const (
	ReceiptKindTerminal ReceiptKind = "terminal"
	ReceiptKindStage    ReceiptKind = "stage"
	ReceiptKindPolicy   ReceiptKind = "policy"
)

type Receipt struct {
	ContractVersion  string               `json:"contract_version"`
	ReceiptID        ids.ReceiptID        `json:"receipt_id"`
	RunID            ids.RunID            `json:"run_id"`
	Kind             ReceiptKind          `json:"kind"`
	EvidenceBundleID ids.EvidenceBundleID `json:"evidence_bundle_id"`
	Summary          string               `json:"summary"`
	ClosureProofRef  string               `json:"closure_proof_ref,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
}
