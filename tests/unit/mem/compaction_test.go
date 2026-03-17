package mem_test

import (
	"testing"

	"srosv2/internal/core/mem"
)

func TestBuildCompactionPlan(t *testing.T) {
	plan := mem.BuildCompactionPlan([]mem.MemoryRecord{{Key: "b"}, {Key: "a"}}, 4)
	if plan.MutationCount != 4 || len(plan.Keys) != 2 || plan.Keys[0] != "a" {
		t.Fatalf("unexpected compaction plan: %+v", plan)
	}
}
