package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHelpQuality(t *testing.T) {
	root := repoRoot(t)
	rootHelp := runHelp(t, root, "--help")
	examplesHelp := runHelp(t, root, "examples", "--help")

	expectedRoot, err := os.ReadFile(filepath.Join(root, "tests", "golden", "cli", "help_root_polished.txt"))
	if err != nil {
		t.Fatalf("read polished root help: %v", err)
	}
	expectedExamples, err := os.ReadFile(filepath.Join(root, "tests", "golden", "cli", "help_examples_polished.txt"))
	if err != nil {
		t.Fatalf("read polished examples help: %v", err)
	}

	if normalizeG(rootHelp) != normalizeG(string(expectedRoot)) {
		t.Fatalf("root help polish mismatch\nEXPECTED:\n%s\nGOT:\n%s", string(expectedRoot), rootHelp)
	}
	if normalizeG(examplesHelp) != normalizeG(string(expectedExamples)) {
		t.Fatalf("examples help polish mismatch\nEXPECTED:\n%s\nGOT:\n%s", string(expectedExamples), examplesHelp)
	}
}

func runHelp(t *testing.T, root string, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", append([]string{"run", "./cmd/sros"}, args...)...)
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("help command failed: %v stderr=%s", err, errOut.String())
	}
	return out.String()
}

func normalizeG(v string) string {
	v = strings.ReplaceAll(v, "\r\n", "\n")
	return strings.TrimSpace(v)
}
