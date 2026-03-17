package cli_test

import (
	"strings"
	"testing"
)

func TestStatusCommand(t *testing.T) {
	code, stdout, stderr := runCLI(t, []string{"status"}, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}
	if !strings.Contains(stdout, "mode: local_cli") {
		t.Fatalf("expected local_cli mode in status output, got:\n%s", stdout)
	}
}
