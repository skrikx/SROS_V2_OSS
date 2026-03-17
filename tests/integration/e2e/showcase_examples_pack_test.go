package e2e_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestShowcaseExamplesPack(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("C:\\Program Files\\Git\\bin\\bash.exe", "./scripts/build_showcase_pack.sh")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("build showcase pack failed: %v stderr=%s", err, errOut.String())
	}
	for _, path := range []string{
		filepath.Join(root, "artifacts", "showcase", "example_catalog.json"),
		filepath.Join(root, "artifacts", "showcase", "first_run_snapshot.json"),
		filepath.Join(root, "artifacts", "showcase", "share_pack_manifest.json"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing showcase artifact %s: %v", path, err)
		}
	}
}
