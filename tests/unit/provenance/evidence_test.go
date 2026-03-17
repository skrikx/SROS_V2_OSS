package provenance_test

import (
	"testing"
	"time"

	"srosv2/contracts/evidence"
	coreprov "srosv2/internal/core/provenance"
)

func TestEmitBundle(t *testing.T) {
	service, err := coreprov.New(t.TempDir(), func() time.Time { return fixedProvNow })
	if err != nil {
		t.Fatalf("new provenance service: %v", err)
	}
	bundle, err := service.EmitBundle("run_001", []evidence.ArtifactRef{{ArtifactID: "art_1", Path: "a", DigestAlgo: evidence.DigestAlgorithmSHA256, Digest: "one"}}, nil, "test")
	if err != nil {
		t.Fatalf("emit bundle: %v", err)
	}
	if bundle.BundleID == "" || len(bundle.HashChain) != 1 {
		t.Fatalf("unexpected bundle: %+v", bundle)
	}
}
