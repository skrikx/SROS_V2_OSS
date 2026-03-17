package cli_test

import (
	"strings"
	"testing"
)

func TestDoctorCommand(t *testing.T) {
	code, stdout, stderr := runCLI(t, []string{"doctor"}, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}
	if !strings.Contains(stdout, "doctor report") {
		t.Fatalf("expected doctor report output, got:\n%s", stdout)
	}
}
