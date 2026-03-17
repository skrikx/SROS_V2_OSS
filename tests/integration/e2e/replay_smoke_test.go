package e2e_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestReplaySmoke(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "replay", "--input", "examples/traces/replay_case_min.json")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("replay failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "replay artifact copied") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
