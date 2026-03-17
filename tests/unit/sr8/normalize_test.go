package sr8_test

import (
	"testing"
	"time"

	"srosv2/internal/core/sr8"
)

func TestNormalize(t *testing.T) {
	parsed := sr8.ParsedIntent{
		Source:      sr8.IntentSourceInline,
		RawIntent:   "  patch   file\n  safely  ",
		OperatorID:  "op_local",
		TenantID:    "local",
		WorkspaceID: "default",
	}
	now := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)
	norm := sr8.Normalize(parsed, now)
	if norm.NormalizedText != "patch file safely" {
		t.Fatalf("unexpected normalized text: %q", norm.NormalizedText)
	}
	if norm.CompileRequestID == "" {
		t.Fatal("compile request id should not be empty")
	}
}
