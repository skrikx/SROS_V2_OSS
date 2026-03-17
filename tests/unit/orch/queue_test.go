package orch_test

import (
	"testing"

	"srosv2/internal/core/orch"
)

func TestQueueFIFO(t *testing.T) {
	queue := orch.NewQueue()
	queue.Enqueue(orch.WorkUnit{ID: "wu-001"})
	queue.Enqueue(orch.WorkUnit{ID: "wu-002"})

	first, err := queue.Dequeue()
	if err != nil {
		t.Fatalf("dequeue first: %v", err)
	}
	second, err := queue.Dequeue()
	if err != nil {
		t.Fatalf("dequeue second: %v", err)
	}
	if first.ID != "wu-001" || second.ID != "wu-002" {
		t.Fatalf("unexpected order: %s then %s", first.ID, second.ID)
	}
}
