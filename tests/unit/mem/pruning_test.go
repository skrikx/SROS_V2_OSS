package mem_test

import (
	"testing"
	"time"

	"srosv2/internal/core/mem"
)

func TestBuildPrunePlan(t *testing.T) {
	plan := mem.BuildPrunePlan([]mem.MemoryRecord{
		{Key: "old", UpdatedAt: fixedMemNow.Add(-48 * time.Hour)},
		{Key: "new", UpdatedAt: fixedMemNow},
	}, fixedMemNow.Add(-24*time.Hour))
	if len(plan.Candidates) != 1 || plan.Candidates[0] != "old" {
		t.Fatalf("unexpected prune plan: %+v", plan)
	}
}
