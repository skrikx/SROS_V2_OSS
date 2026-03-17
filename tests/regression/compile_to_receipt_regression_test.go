package regression_test

import (
	"os"
	"testing"
)

func TestCompileToReceiptFixturesExist(t *testing.T) {
	for _, path := range []string{
		"../../tests/golden/compile/compile_receipt_min.json",
		"../../tests/golden/provenance/receipt_min.json",
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing fixture %s: %v", path, err)
		}
	}
}
