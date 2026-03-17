package evidence

import "srosv2/internal/shared/ids"

type Bundle struct {
	BundleID     ids.EvidenceBundleID `json:"bundle_id"`
	RunID        ids.RunID            `json:"run_id"`
	ArtifactRefs []ArtifactRef        `json:"artifact_refs"`
	ReceiptRefs  []ids.ReceiptID      `json:"receipt_refs,omitempty"`
	HashChain    []string             `json:"hash_chain,omitempty"`
	Notes        string               `json:"notes,omitempty"`
}
