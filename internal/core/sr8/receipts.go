package sr8

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/contracts/evidence"
	"srosv2/contracts/runcontract"
	"srosv2/internal/shared/ids"
)

func BuildCompileReceipt(
	requestID string,
	class Classification,
	topo TopologyDecision,
	contract runcontract.RunContract,
	bundleArtifacts []evidence.ArtifactRef,
	now time.Time,
) (CompileReceipt, error) {
	receipt := evidence.Receipt{
		ContractVersion:  "v2.0",
		ReceiptID:        ids.ReceiptID("receipt_" + shortHash(requestID+"|receipt")),
		RunID:            contract.RunID,
		Kind:             evidence.ReceiptKindStage,
		EvidenceBundleID: ids.EvidenceBundleID("bundle_" + shortHash(requestID+"|bundle")),
		Status:           "compiled",
		ArtifactRefs:     bundleArtifacts,
		Summary:          "compile completed; runtime admission not performed",
		ClosureProofRef:  "sr8-compile-only",
		CreatedAt:        now.UTC(),
	}
	bundle := evidence.Bundle{
		BundleID:     receipt.EvidenceBundleID,
		RunID:        contract.RunID,
		ArtifactRefs: bundleArtifacts,
		Notes:        "compile plane evidence only; trace linkage provisional",
	}

	if errs := evidence.ValidateReceipt(receipt); len(errs) > 0 {
		return CompileReceipt{}, fmt.Errorf("invalid receipt: %v", errs[0])
	}
	if errs := evidence.ValidateBundle(bundle); len(errs) > 0 {
		return CompileReceipt{}, fmt.Errorf("invalid evidence bundle: %v", errs[0])
	}

	return CompileReceipt{
		Receipt:          receipt,
		Bundle:           bundle,
		CompileRequestID: requestID,
		DomainClass:      class.Domain,
		RiskClass:        class.Risk,
		TopologyClass:    topo.Topology,
		Status:           "compiled",
		TraceLinkage:     "provisional: trace append deferred to W08",
		RuntimeAdmission: "not_admitted",
	}, nil
}

func ArtifactRefFromFile(path string) (evidence.ArtifactRef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return evidence.ArtifactRef{}, err
	}
	digest := sha256.Sum256(data)

	return evidence.ArtifactRef{
		ArtifactID: ids.ArtifactID("artifact_" + shortHash(path)),
		Path:       filepath.Clean(path),
		DigestAlgo: evidence.DigestAlgorithmSHA256,
		Digest:     hex.EncodeToString(digest[:]),
		MediaType:  mediaTypeForPath(path),
		Bytes:      int64(len(data)),
	}, nil
}

func mediaTypeForPath(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		return "application/json"
	case ".srxml", ".xml":
		return "application/xml"
	default:
		return "application/octet-stream"
	}
}
