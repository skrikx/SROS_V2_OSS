package cli_integration_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestToolsSearchSmoke(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "tools", "search", "--class", "tool.local")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("tools search failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "tool search completed") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
