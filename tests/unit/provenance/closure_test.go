package provenance_test

import (
	"testing"
	"time"

	coreprov "srosv2/internal/core/provenance"
)

func TestEmitClosure(t *testing.T) {
	service, err := coreprov.New(t.TempDir(), func() time.Time { return fixedProvNow })
	if err != nil {
		t.Fatalf("new provenance service: %v", err)
	}
	proof, ref, err := service.EmitClosure("run_001", "rolled_back", []string{"evt_1"}, []string{"rcpt_1"}, []string{"art_1"})
	if err != nil {
		t.Fatalf("emit closure: %v", err)
	}
	if proof.ClosureStatus != "sealed" || ref == "" {
		t.Fatalf("unexpected closure proof: %+v ref=%s", proof, ref)
	}
}
