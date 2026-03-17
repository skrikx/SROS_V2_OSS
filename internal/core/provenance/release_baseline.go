package provenance

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/contracts/release"
	"srosv2/internal/shared/ids"
)

type ReleaseBaseline struct {
	root  string
	now   func() time.Time
	store *PostgresStore
}

func NewReleaseBaseline(root string, now func() time.Time, store *PostgresStore) (*ReleaseBaseline, error) {
	if root == "" {
		root = filepath.Join("artifacts", "releases")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create release baseline root: %w", err)
	}
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &ReleaseBaseline{root: root, now: now, store: store}, nil
}

func (r *ReleaseBaseline) Pack(ctx context.Context, checkpointID ids.CheckpointID, stage release.Stage, meta map[string]string) (map[string]any, error) {
	record := release.ReleaseRecord{
		ContractVersion: "v2.0",
		ReleaseID:       ids.ReleaseID("rel_" + digestBytes([]byte(string(checkpointID) + string(stage)))[:12]),
		CheckpointID:    checkpointID,
		TargetStage:     stage,
		PromotionMeta:   meta,
		CreatedAt:       r.now().UTC(),
	}
	if errs := release.ValidateRelease(record); len(errs) > 0 {
		return nil, errs[0]
	}
	gateID := "gate_" + string(record.ReleaseID)
	gate := map[string]any{
		"gate_id":    gateID,
		"release_id": record.ReleaseID,
		"status":     "pass",
		"checked_at": r.now().UTC(),
	}
	data, err := json.MarshalIndent(map[string]any{"release": record, "promotion_gate": gate}, "", "  ")
	if err != nil {
		return nil, err
	}
	out := filepath.Join(r.root, string(record.ReleaseID)+".json")
	if err := os.WriteFile(out, append(data, '\n'), 0o644); err != nil {
		return nil, err
	}
	if r.store != nil {
		if err := r.store.SaveRelease(ctx, record); err != nil {
			return nil, err
		}
		if err := r.store.SavePromotionGateResult(ctx, gateID, string(record.ReleaseID), "pass", gate); err != nil {
			return nil, err
		}
	}
	return map[string]any{"release": record, "promotion_gate": gate, "artifact": out}, nil
}
