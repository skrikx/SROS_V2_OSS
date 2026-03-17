package cli_integration_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	cmemory "srosv2/contracts/memory"
)

func TestMemoryBranchRewindSmoke(t *testing.T) {
	root := repoRoot(t)
	seed := filepath.Join(root, "examples", "memory", "workspace_seed.json")
	branch := filepath.Join(root, "examples", "memory", "branch_lineage.json")
	runCLICommand(t, root, "memory", "recall", "--input", seed)
	out := runCLICommand(t, root, "memory", "branch", "--input", branch)
	if !strings.Contains(out, "memory branch applied") {
		t.Fatalf("unexpected branch output: %s", out)
	}

	branchData, err := os.ReadFile(filepath.Join(root, "artifacts", "memory", "branches", "branch_feature.json"))
	if err != nil {
		t.Fatalf("read branch record: %v", err)
	}
	var record cmemory.BranchRecord
	if err := json.Unmarshal(branchData, &record); err != nil {
		t.Fatalf("decode branch record: %v", err)
	}
	rewindOut := runCLICommand(t, root, "memory", "rewind", "--branch", "branch_feature", "--mutation", string(record.HeadMutationID), "--operator", "op_local")
	if !strings.Contains(rewindOut, "memory rewind applied") {
		t.Fatalf("unexpected rewind output: %s", rewindOut)
	}
}

func runCLICommand(t *testing.T, root string, args ...string) string {
	t.Helper()
	var out, errBuf bytes.Buffer
	cmd := exec.Command("go", append([]string{"run", "./cmd/sros"}, args...)...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("cli command failed: %v stderr=%s", err, errBuf.String())
	}
	return out.String()
}
