package provenance_test

import (
	"testing"
	"time"

	"srosv2/contracts/evidence"
	coreprov "srosv2/internal/core/provenance"
)

func TestEmitReceipt(t *testing.T) {
	service, err := coreprov.New(t.TempDir(), func() time.Time { return fixedProvNow })
	if err != nil {
		t.Fatalf("new provenance service: %v", err)
	}
	receipt, err := service.EmitReceipt("run_001", evidence.ReceiptKindTerminal, "sealed", "run complete", nil, "")
	if err == nil {
		t.Fatal("expected receipt emission without artifacts to fail")
	}
	receipt, err = service.EmitReceipt("run_001", evidence.ReceiptKindTerminal, "sealed", "run complete", []evidence.ArtifactRef{
		{ArtifactID: "art_1", Path: "artifact.json", DigestAlgo: evidence.DigestAlgorithmSHA256, Digest: "abc"},
	}, "")
	if err != nil {
		t.Fatalf("emit receipt: %v", err)
	}
	if receipt.ReceiptID == "" || receipt.Status != "sealed" {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
}
