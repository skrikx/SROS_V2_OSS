package mirror_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"srosv2/internal/core/mirror"
)

func TestWitnessWritten(t *testing.T) {
	root := t.TempDir()
	engine, err := mirror.New(root, func() time.Time { return fixedMirrorNow })
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	event, _, err := engine.Observe(mirror.RuntimeSnapshot{RunID: "run_001", RuntimeState: "running", MemoryMutations: 1}, "test")
	if err != nil {
		t.Fatalf("observe: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "witness", event.WitnessID+".json"))
	if err != nil {
		t.Fatalf("read witness: %v", err)
	}
	var decoded mirror.WitnessEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode witness: %v", err)
	}
}
