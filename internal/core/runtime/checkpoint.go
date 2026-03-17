package runtime

import (
	"fmt"
	"strings"
	"time"

	"srosv2/contracts/release"
	"srosv2/internal/shared/ids"
)

type RuntimeCheckpoint struct {
	SessionID string                   `json:"session_id"`
	State     SessionState             `json:"state"`
	Record    release.CheckpointRecord `json:"record"`
}

func NewCheckpoint(session RuntimeSession, stage string, now time.Time) (RuntimeCheckpoint, error) {
	releaseStage, err := parseCheckpointStage(stage)
	if err != nil {
		return RuntimeCheckpoint{}, err
	}

	record := release.CheckpointRecord{
		ContractVersion: "v2.0",
		CheckpointID:    ids.CheckpointID("cp_" + shortHash(session.SessionID+"|"+now.UTC().Format(time.RFC3339Nano))),
		RunID:           ids.RunID(session.RunID),
		Stage:           releaseStage,
		RecordedAt:      now.UTC(),
	}
	if errs := release.ValidateCheckpoint(record); len(errs) > 0 {
		return RuntimeCheckpoint{}, fmt.Errorf("invalid checkpoint record: %v", errs[0])
	}

	return RuntimeCheckpoint{
		SessionID: session.SessionID,
		State:     session.State,
		Record:    record,
	}, nil
}

func parseCheckpointStage(stage string) (release.Stage, error) {
	s := strings.TrimSpace(strings.ToLower(stage))
	if s == "" {
		return release.StageValidated, nil
	}
	switch release.Stage(s) {
	case release.StageDraft, release.StageValidated, release.StagePromoted, release.StageRolledBack:
		return release.Stage(s), nil
	default:
		return "", fmt.Errorf("unsupported checkpoint stage %q", stage)
	}
}
