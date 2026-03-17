package cli_integration_test

import (
	"strings"
	"testing"
)

func TestReceiptExportSmoke(t *testing.T) {
	root := repoRoot(t)
	out := runTraceCLICommand(t, root, "receipts", "export", "--input", joinRoot(root, "examples", "provenance", "receipt_bundle_min.json"))
	if !strings.Contains(out, "receipt bundle exported") {
		t.Fatalf("unexpected output: %s", out)
	}
}
