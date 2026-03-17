package cli_test

import (
	"strings"
	"testing"
)

func TestInspectCommand(t *testing.T) {
	code, stdout, stderr := runCLI(t, []string{"inspect"}, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}
	if !strings.Contains(stdout, "repository wiring snapshot") {
		t.Fatalf("expected inspect snapshot output, got:\n%s", stdout)
	}
}
