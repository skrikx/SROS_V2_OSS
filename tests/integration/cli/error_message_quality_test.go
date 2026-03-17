package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestErrorMessageQuality(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "receipts", "export")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected command to fail")
	}
	expected, readErr := os.ReadFile(filepath.Join(root, "tests", "golden", "cli", "error_missing_input_polished.txt"))
	if readErr != nil {
		t.Fatalf("read error golden: %v", readErr)
	}
	got := normalizeG(errOut.String())
	if !strings.Contains(got, normalizeG(string(expected))) {
		t.Fatalf("unexpected polished error\nEXPECTED CONTAINS:\n%s\nGOT:\n%s", string(expected), errOut.String())
	}
}
