package cli_integration_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestConnectorEnvelopeSmoke(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "connectors", "envelope", "inspect", "--input", "examples/connectors/local_secret_envelope.json")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("connector envelope inspect failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "connector envelope inspected") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
