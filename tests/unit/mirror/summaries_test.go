package mirror_test

import (
	"testing"

	"srosv2/internal/core/mirror"
)

func TestBuildSummary(t *testing.T) {
	summary := mirror.BuildSummary(mirror.RuntimeSnapshot{RunID: "run_001"}, mirror.DriftFlag{Level: "low"}, 1)
	if summary.SourceBasis == "" {
		t.Fatal("expected source basis")
	}
}
