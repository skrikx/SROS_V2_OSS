package cli_integration_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestToolsValidateSmoke(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "tools", "validate", "--manifest", "examples/tools/local_patch_manifest.json")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("tools validate failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "tool manifest validated") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
