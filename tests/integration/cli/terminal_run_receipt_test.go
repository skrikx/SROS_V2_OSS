package cli_integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTerminalRunReceipt(t *testing.T) {
	root := repoRoot(t)
	runTraceCLICommand(t, root, "run", "--contract", joinRoot(root, "examples", "runs", "minimal_run_contract.json"))
	runTraceCLICommand(t, root, "checkpoint", "--latest")
	runTraceCLICommand(t, root, "rollback")
	closures, err := os.ReadDir(filepath.Join(root, "artifacts", "provenance", "closures"))
	if err != nil || len(closures) == 0 {
		t.Fatalf("expected closure proofs: %v", err)
	}
	out := runTraceCLICommand(t, root, "receipts", "closure", "--input", filepath.Join(root, "artifacts", "provenance", "closures", closures[0].Name()))
	if !strings.Contains(out, "closure proof inspected") {
		t.Fatalf("unexpected output: %s", out)
	}
}
