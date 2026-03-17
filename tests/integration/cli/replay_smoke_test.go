package cli_integration_test

import (
	"strings"
	"testing"
)

func TestReplaySmoke(t *testing.T) {
	root := repoRoot(t)
	runTraceCLICommand(t, root, "run", "--contract", joinRoot(root, "examples", "runs", "minimal_run_contract.json"))
	out := runTraceCLICommand(t, root, "trace", "replay", "--run-id", "run_minimal_001")
	if !strings.Contains(out, "trace replay completed") {
		t.Fatalf("unexpected output: %s", out)
	}
}
