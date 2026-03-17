package mirror_test

import (
	"testing"
	"time"

	"srosv2/internal/core/mirror"
)

func TestEngineObserve(t *testing.T) {
	engine, err := mirror.New(t.TempDir(), func() time.Time { return fixedMirrorNow })
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	event, summary, err := engine.Observe(mirror.RuntimeSnapshot{
		RunID:           "run_001",
		SessionID:       "sess_001",
		RuntimeState:    "running",
		MemoryMutations: 1,
		BranchCount:     1,
	}, "test")
	if err != nil {
		t.Fatalf("observe: %v", err)
	}
	if event.WitnessID == "" || summary.DriftLevel == "" {
		t.Fatalf("unexpected mirror output: %+v %+v", event, summary)
	}
}
