package provenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"srosv2/contracts/evidence"
	"srosv2/internal/shared/ids"
)

func (s *Service) EmitReceipt(runID ids.RunID, kind evidence.ReceiptKind, status, summary string, artifacts []evidence.ArtifactRef, closureRef string) (evidence.Receipt, error) {
	bundle, err := s.EmitBundle(runID, artifacts, nil, summary)
	if err != nil {
		return evidence.Receipt{}, err
	}
	receipt := evidence.Receipt{
		ContractVersion:  "v2.0",
		ReceiptID:        ids.ReceiptID("rcpt_" + digestBytes([]byte(string(runID) + summary))[:12]),
		RunID:            runID,
		Kind:             kind,
		EvidenceBundleID: bundle.BundleID,
		Status:           status,
		ArtifactRefs:     artifacts,
		Summary:          summary,
		ClosureProofRef:  closureRef,
		CreatedAt:        s.now().UTC(),
	}
	if errs := evidence.ValidateReceipt(receipt); len(errs) > 0 {
		return evidence.Receipt{}, errs[0]
	}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return evidence.Receipt{}, fmt.Errorf("marshal receipt: %w", err)
	}
	path := filepath.Join(s.root, "receipts", string(receipt.ReceiptID)+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return evidence.Receipt{}, fmt.Errorf("write receipt: %w", err)
	}
	return receipt, nil
}
