package shared_test

import (
	"testing"

	"srosv2/internal/shared/ids"
)

func TestParseRunIDValid(t *testing.T) {
	runID, err := ids.ParseRunID("run_001")
	if err != nil {
		t.Fatalf("expected valid run id: %v", err)
	}
	if runID != "run_001" {
		t.Fatalf("unexpected run id: %s", runID)
	}
}

func TestParseRunIDInvalid(t *testing.T) {
	if _, err := ids.ParseRunID("bad id"); err == nil {
		t.Fatal("expected parse error for unsafe id")
	}
}
