package e2e_test

import (
	"bytes"
	"strings"
	"testing"
)

func TestFirstRunTimeToValue(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := bashScriptCommand(t, "./scripts/first_run_smoke.sh")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("first run smoke failed: %v stderr=%s", err, errOut.String())
	}
	for _, token := range []string{"verification completed", "examples catalog loaded", "trace inspected", "receipt bundle exported", "tool search completed"} {
		if !strings.Contains(out.String(), token) {
			t.Fatalf("missing token %q in first run output:\n%s", token, out.String())
		}
	}
}
