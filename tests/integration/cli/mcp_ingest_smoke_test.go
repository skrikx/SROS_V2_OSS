package cli_integration_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestMCPIngestSmoke(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "mcp", "ingest", "--input", "examples/mcp/ingested_remote_capability.json")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("mcp ingest failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "normalized and admitted") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
