package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMemoryRecallSmoke(t *testing.T) {
	root := repoRoot(t)
	input := filepath.Join(root, "examples", "memory", "workspace_seed.json")
	var out, errBuf bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "memory", "recall", "--input", input)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("memory recall failed: %v stderr=%s", err, errBuf.String())
	}
	if !strings.Contains(out.String(), "memory recall completed") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
