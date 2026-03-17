package provenance

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"srosv2/contracts/evidence"
	"srosv2/contracts/release"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) SaveReceipt(ctx context.Context, receipt evidence.Receipt) error {
	data, err := json.Marshal(receipt)
	if err != nil {
		return fmt.Errorf("marshal receipt: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO receipts (receipt_id, run_id, kind, receipt_json, created_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (receipt_id) DO UPDATE SET receipt_json = EXCLUDED.receipt_json`,
		receipt.ReceiptID, receipt.RunID, receipt.Kind, data, receipt.CreatedAt)
	return err
}

func (s *PostgresStore) SaveBundle(ctx context.Context, bundle evidence.Bundle) error {
	data, err := json.Marshal(bundle)
	if err != nil {
		return fmt.Errorf("marshal evidence bundle: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO evidence_bundles (bundle_id, run_id, bundle_json, created_at)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (bundle_id) DO UPDATE SET bundle_json = EXCLUDED.bundle_json`,
		bundle.BundleID, bundle.RunID, data, time.Now().UTC())
	return err
}

func (s *PostgresStore) SaveArtifactProvenance(ctx context.Context, prov ArtifactProvenance) error {
	data, err := json.Marshal(prov)
	if err != nil {
		return fmt.Errorf("marshal artifact provenance: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO artifact_provenance (artifact_id, run_id, source_kind, artifact_json, linked_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (artifact_id) DO UPDATE SET artifact_json = EXCLUDED.artifact_json`,
		prov.Artifact.ArtifactID, prov.RunID, prov.SourceKind, data, prov.LinkedAt)
	return err
}

func (s *PostgresStore) SaveRelease(ctx context.Context, record release.ReleaseRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal release record: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO releases (release_id, checkpoint_id, target_stage, release_json, created_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (release_id) DO UPDATE SET release_json = EXCLUDED.release_json`,
		record.ReleaseID, record.CheckpointID, record.TargetStage, data, record.CreatedAt)
	return err
}

func (s *PostgresStore) SavePromotionGateResult(ctx context.Context, gateID, releaseID, status string, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal promotion gate result: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO promotion_gate_results (gate_id, release_id, status, result_json, created_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (gate_id) DO UPDATE SET result_json = EXCLUDED.result_json, status = EXCLUDED.status`,
		gateID, releaseID, status, data, time.Now().UTC())
	return err
}

func (s *PostgresStore) SaveRollback(ctx context.Context, record release.RollbackRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal release rollback: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO rollback_records (rollback_id, release_id, target_checkpoint_id, rollback_json, created_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (rollback_id) DO UPDATE SET rollback_json = EXCLUDED.rollback_json`,
		record.RollbackID, record.ReleaseID, record.TargetCheckpointID, data, record.CreatedAt)
	return err
}
