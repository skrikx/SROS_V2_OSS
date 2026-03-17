package contracts_test

import (
	"encoding/json"
	"testing"
	"time"

	"srosv2/contracts/release"
	"srosv2/internal/shared/ids"
)

func TestCheckpointValidateValid(t *testing.T) {
	record := release.CheckpointRecord{
		ContractVersion: "v2.0",
		CheckpointID:    ids.CheckpointID("cp_001"),
		RunID:           ids.RunID("run_001"),
		Stage:           release.StageValidated,
		RecordedAt:      time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
	}

	errs := release.ValidateCheckpoint(record)
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(errs))
	}
}

func TestReleaseGoldenFixtures(t *testing.T) {
	checkpointData := loadFixture(t, "checkpoint_record.json")
	releaseData := loadFixture(t, "release_record.json")

	var checkpoint release.CheckpointRecord
	if err := json.Unmarshal(checkpointData, &checkpoint); err != nil {
		t.Fatalf("unmarshal checkpoint fixture: %v", err)
	}
	if errs := release.ValidateCheckpoint(checkpoint); len(errs) != 0 {
		t.Fatalf("checkpoint fixture invalid: %d errors", len(errs))
	}

	var record release.ReleaseRecord
	if err := json.Unmarshal(releaseData, &record); err != nil {
		t.Fatalf("unmarshal release fixture: %v", err)
	}
	if errs := release.ValidateRelease(record); len(errs) != 0 {
		t.Fatalf("release fixture invalid: %d errors", len(errs))
	}
}
