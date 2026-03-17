package mem_test

import (
	"testing"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/core/mem"
)

func TestBuildSessionTree(t *testing.T) {
	tree := mem.BuildSessionTree([]cmemory.MemoryMutation{
		{RunID: "run_001", SessionID: "sess_001", Key: "a"},
		{RunID: "run_001", SessionID: "sess_001", Key: "b"},
		{RunID: "run_002", SessionID: "sess_002", Key: "c"},
	})
	if len(tree) != 2 {
		t.Fatalf("expected 2 session nodes, got %d", len(tree))
	}
}
