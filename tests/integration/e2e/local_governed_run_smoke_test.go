package e2e_test

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestLocalGovernedRunSmoke(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "verify")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("verify failed: %v stderr=%s", err, errOut.String())
	}
}
