package provenance_test

import (
	"testing"

	"srosv2/contracts/evidence"
	coreprov "srosv2/internal/core/provenance"
)

func TestHashChain(t *testing.T) {
	chain, err := coreprov.HashChain([]evidence.ArtifactRef{
		{ArtifactID: "art_1", Path: "a", DigestAlgo: evidence.DigestAlgorithmSHA256, Digest: "one"},
		{ArtifactID: "art_2", Path: "b", DigestAlgo: evidence.DigestAlgorithmSHA256, Digest: "two"},
	})
	if err != nil {
		t.Fatalf("hash chain: %v", err)
	}
	if len(chain) != 2 {
		t.Fatalf("expected 2 hashes, got %d", len(chain))
	}
}
