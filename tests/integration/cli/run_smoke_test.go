package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSmoke(t *testing.T) {
	root := repoRoot(t)
	contract := filepath.Join(root, "examples", "runs", "minimal_run_contract.json")

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "run", "--contract", contract)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("run smoke failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "runtime session admitted") {
		t.Fatalf("unexpected run output: %s", out.String())
	}
}
