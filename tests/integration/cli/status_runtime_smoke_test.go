package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatusRuntimeSmoke(t *testing.T) {
	root := repoRoot(t)
	contract := filepath.Join(root, "examples", "runs", "minimal_run_contract.json")

	runCLICmd(t, root, "run", "--contract", contract)

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "status", "--latest")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("status smoke failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "runtime_state:") {
		t.Fatalf("unexpected status output: %s", out.String())
	}
}
