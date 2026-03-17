package cli_test

import (
	"strings"
	"testing"
)

func TestRootHelpGolden(t *testing.T) {
	code, stdout, stderr := runCLI(t, []string{"--help"}, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}

	expected := loadGolden(t, "help_root.txt")
	if normalize(stdout) != normalize(expected) {
		t.Fatalf("root help mismatch\nEXPECTED:\n%s\nGOT:\n%s", expected, stdout)
	}
}

func TestRunHelpGolden(t *testing.T) {
	code, stdout, stderr := runCLI(t, []string{"run", "--help"}, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}
	expected := loadGolden(t, "help_run.txt")
	if normalize(stdout) != normalize(expected) {
		t.Fatalf("run help mismatch\nEXPECTED:\n%s\nGOT:\n%s", expected, stdout)
	}
}

func TestToolsHelpGolden(t *testing.T) {
	code, stdout, stderr := runCLI(t, []string{"tools", "--help"}, nil)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}
	expected := loadGolden(t, "help_tools.txt")
	if normalize(stdout) != normalize(expected) {
		t.Fatalf("tools help mismatch\nEXPECTED:\n%s\nGOT:\n%s", expected, stdout)
	}
}

func normalize(v string) string {
	v = strings.ReplaceAll(v, "\r\n", "\n")
	return strings.TrimSpace(v)
}
