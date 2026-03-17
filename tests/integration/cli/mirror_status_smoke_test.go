package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMirrorStatusSmoke(t *testing.T) {
	root := repoRoot(t)
	input := filepath.Join(root, "examples", "mirror", "runtime_snapshot.json")
	var out, errBuf bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "mirror", "status", "--input", input)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("mirror status failed: %v stderr=%s", err, errBuf.String())
	}
	if !strings.Contains(out.String(), "mirror status captured") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
