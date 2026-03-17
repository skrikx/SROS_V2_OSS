package cli_integration_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestHelpSmoke(t *testing.T) {
	root := repoRoot(t)

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "--help")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("root help failed: %v stderr=%s", err, errOut.String())
	}

	out.Reset()
	errOut.Reset()
	cmd = exec.Command("go", "run", "./cmd/sros", "run", "--help")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("run help failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "Admit a canonical run contract into SR9 runtime") {
		t.Fatalf("unexpected run help output: %s", out.String())
	}
}
