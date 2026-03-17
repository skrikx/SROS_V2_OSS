package sr8_test

import (
	"os"
	"path/filepath"
	"testing"

	"srosv2/internal/core/sr8"
)

func TestParseRequestInline(t *testing.T) {
	parsed, err := sr8.ParseRequest(sr8.CompileRequest{Intent: "hello world"})
	if err != nil {
		t.Fatalf("parse inline request: %v", err)
	}
	if parsed.Source != sr8.IntentSourceInline {
		t.Fatalf("expected inline source, got %s", parsed.Source)
	}
}

func TestParseRequestFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "intent.txt")
	if err := os.WriteFile(input, []byte("patch this file"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	parsed, err := sr8.ParseRequest(sr8.CompileRequest{InputPath: input})
	if err != nil {
		t.Fatalf("parse file request: %v", err)
	}
	if parsed.Source != sr8.IntentSourceFile {
		t.Fatalf("expected file source, got %s", parsed.Source)
	}
}
