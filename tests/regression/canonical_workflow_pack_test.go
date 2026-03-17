package regression_test

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCanonicalWorkflowPackExists(t *testing.T) {
	data, err := os.ReadFile("../../examples/regression/canonical_workflow_pack.json")
	if err != nil {
		t.Fatalf("read workflow pack: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("decode workflow pack: %v", err)
	}
	if payload["name"] != "canonical_workflow_pack" {
		t.Fatalf("unexpected workflow pack name")
	}
}
