package provenance_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	coreprov "srosv2/internal/core/provenance"
)

func TestLinkArtifact(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "artifact.json")
	if err := os.WriteFile(file, []byte("{\"ok\":true}\n"), 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	service, err := coreprov.New(root, func() time.Time { return fixedProvNow })
	if err != nil {
		t.Fatalf("new provenance service: %v", err)
	}
	ref, err := service.LinkArtifact("run_001", file, "application/json", "test")
	if err != nil {
		t.Fatalf("link artifact: %v", err)
	}
	if ref.ArtifactID == "" {
		t.Fatal("expected artifact id")
	}
}
