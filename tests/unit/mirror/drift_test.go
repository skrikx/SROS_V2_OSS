package mirror_test

import (
	"testing"

	"srosv2/internal/core/mirror"
)

func TestDetectDriftHigh(t *testing.T) {
	flag := mirror.DetectDrift(mirror.RuntimeSnapshot{RuntimeState: "failed_safe"})
	if flag.Level != "high" {
		t.Fatalf("expected high drift, got %+v", flag)
	}
}
