package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBootstrapSmoke(t *testing.T) {
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "bootstrap")
	cmd.Dir = repoRoot(t)
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("bootstrap failed: %v stderr=%s", err, errBuf.String())
	}
	if out.Len() == 0 {
		t.Fatal("expected bootstrap output")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", "..", ".."))
}
