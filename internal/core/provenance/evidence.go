package provenance

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"srosv2/contracts/evidence"
	"srosv2/internal/shared/ids"
)

func (s *Service) EmitBundle(runID ids.RunID, artifacts []evidence.ArtifactRef, receiptRefs []ids.ReceiptID, notes string) (evidence.Bundle, error) {
	chain, err := HashChain(artifacts)
	if err != nil {
		return evidence.Bundle{}, err
	}
	bundle := evidence.Bundle{
		BundleID:     ids.EvidenceBundleID("bundle_" + digestBytes([]byte(string(runID) + notes))[:12]),
		RunID:        runID,
		ArtifactRefs: artifacts,
		ReceiptRefs:  receiptRefs,
		HashChain:    chain,
		Notes:        notes,
	}
	if errs := evidence.ValidateBundle(bundle); len(errs) > 0 {
		return evidence.Bundle{}, errs[0]
	}
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return evidence.Bundle{}, fmt.Errorf("marshal bundle: %w", err)
	}
	path := filepath.Join(s.root, "bundles", string(bundle.BundleID)+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return evidence.Bundle{}, fmt.Errorf("write bundle: %w", err)
	}
	if s.pgStore != nil {
		if err := s.pgStore.SaveBundle(context.Background(), bundle); err != nil {
			return evidence.Bundle{}, err
		}
	}
	return bundle, nil
}
