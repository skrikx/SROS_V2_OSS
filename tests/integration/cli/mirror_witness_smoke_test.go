package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMirrorWitnessSmoke(t *testing.T) {
	root := repoRoot(t)
	input := filepath.Join(root, "examples", "mirror", "witness_case.json")
	var out, errBuf bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "mirror", "witness", "--input", input)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("mirror witness failed: %v stderr=%s", err, errBuf.String())
	}
	if !strings.Contains(out.String(), "mirror witness emitted") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
