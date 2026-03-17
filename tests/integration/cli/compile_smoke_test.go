package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompileSmoke(t *testing.T) {
	root := repoRoot(t)
	input := filepath.Join(root, "examples", "intents", "minimal.txt")

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "compile", "--input", input)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("compile smoke failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "compile accepted") {
		t.Fatalf("unexpected compile output: %s", out.String())
	}
}
