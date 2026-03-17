package contracts_test

import (
	"encoding/json"
	"testing"
	"time"

	"srosv2/contracts/evidence"
	"srosv2/internal/shared/ids"
)

func TestReceiptValidateValid(t *testing.T) {
	receipt := evidence.Receipt{
		ContractVersion:  "v2.0",
		ReceiptID:        ids.ReceiptID("receipt_001"),
		RunID:            ids.RunID("run_001"),
		Kind:             evidence.ReceiptKindTerminal,
		EvidenceBundleID: ids.EvidenceBundleID("bundle_001"),
		Status:           "sealed",
		Summary:          "run complete",
		CreatedAt:        time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
	}

	errs := evidence.ValidateReceipt(receipt)
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(errs))
	}
}

func TestReceiptGoldenFixture(t *testing.T) {
	data := loadFixture(t, "receipt.json")
	var receipt evidence.Receipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	errs := evidence.ValidateReceipt(receipt)
	if len(errs) != 0 {
		t.Fatalf("expected fixture to validate, got %d errors", len(errs))
	}
}
