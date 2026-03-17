package sr8_test

import (
	"testing"

	"srosv2/internal/core/sr8"
)

func TestValidateCompileInput(t *testing.T) {
	err := sr8.ValidateCompileInput(sr8.NormalizedIntent{}, sr8.Classification{}, sr8.TopologyDecision{})
	if err == nil {
		t.Fatal("expected compile input validation error")
	}
}
