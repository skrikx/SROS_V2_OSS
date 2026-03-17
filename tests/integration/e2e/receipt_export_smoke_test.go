package e2e_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestReceiptExportSmoke(t *testing.T) {
	root := repoRoot(t)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := exec.Command("go", "run", "./cmd/sros", "receipts", "export", "--input", "examples/provenance/receipt_bundle_min.json")
	cmd.Dir = root
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("receipt export failed: %v stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "receipt bundle exported") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
