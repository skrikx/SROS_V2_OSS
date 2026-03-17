package eval_test

import (
	"os"
	"path/filepath"
	"testing"

	"srosv2/internal/core/provenance"
	"srosv2/internal/shared/ids"
)

func TestEvidenceBindingArtifactsRemainLinkable(t *testing.T) {
	root := t.TempDir()
	service, err := provenance.New(filepath.Join(root, "prov"), nil)
	if err != nil {
		t.Fatalf("new provenance: %v", err)
	}
	artifact := filepath.Join(root, "artifact.json")
	if err := os.WriteFile(artifact, []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	ref, err := service.LinkArtifact(ids.RunID("run_eval"), artifact, "application/json", "eval")
	if err != nil {
		t.Fatalf("link artifact: %v", err)
	}
	if ref.Path == "" {
		t.Fatalf("expected artifact path")
	}
}
