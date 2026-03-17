package runtime

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) SaveSession(ctx context.Context, session RuntimeSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal runtime session: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO run_contracts (run_id, tenant_id, workspace_id, contract_json, created_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (run_id) DO UPDATE SET contract_json = EXCLUDED.contract_json`, session.RunID, session.Contract.TenantID, session.Contract.WorkspaceID, dataJSON(session.Contract), session.CreatedAt); err != nil {
		return fmt.Errorf("upsert run contract: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO runs (session_id, run_id, tenant_id, workspace_id, state, session_json, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (session_id) DO UPDATE SET state = EXCLUDED.state, session_json = EXCLUDED.session_json, updated_at = EXCLUDED.updated_at`,
		session.SessionID, session.RunID, session.Contract.TenantID, session.Contract.WorkspaceID, session.State, data, session.CreatedAt, session.UpdatedAt); err != nil {
		return fmt.Errorf("upsert runtime session: %w", err)
	}
	for _, item := range session.History {
		if _, err := s.db.ExecContext(ctx, `
			INSERT INTO run_state_transitions (session_id, run_id, from_state, to_state, reason, occurred_at)
			VALUES ($1,$2,$3,$4,$5,$6)`,
			session.SessionID, session.RunID, item.From, item.To, item.Reason, item.At); err != nil {
			return fmt.Errorf("insert runtime transition: %w", err)
		}
	}
	return nil
}

func (s *PostgresStore) SaveCheckpoint(ctx context.Context, cp RuntimeCheckpoint) error {
	data, err := json.Marshal(cp)
	if err != nil {
		return fmt.Errorf("marshal checkpoint: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO checkpoints (checkpoint_id, session_id, run_id, stage, checkpoint_json, recorded_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (checkpoint_id) DO UPDATE SET checkpoint_json = EXCLUDED.checkpoint_json`,
		cp.Record.CheckpointID, cp.SessionID, cp.Record.RunID, cp.Record.Stage, data, cp.Record.RecordedAt)
	return err
}

func (s *PostgresStore) SaveRollback(ctx context.Context, rb RuntimeRollback) error {
	data, err := json.Marshal(rb)
	if err != nil {
		return fmt.Errorf("marshal rollback: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO rollback_records (rollback_id, release_id, target_checkpoint_id, rollback_json, created_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (rollback_id) DO UPDATE SET rollback_json = EXCLUDED.rollback_json`,
		rb.Record.RollbackID, rb.Record.ReleaseID, rb.Record.TargetCheckpointID, data, rb.Record.CreatedAt)
	return err
}

func (s *PostgresStore) SaveApproval(ctx context.Context, session RuntimeSession, approval ApprovalCheckpoint) error {
	requestedAt, _ := time.Parse(time.RFC3339, approval.RequestedAt)
	var approvedAt any
	if approval.ApprovedAt != "" {
		if t, err := time.Parse(time.RFC3339, approval.ApprovedAt); err == nil {
			approvedAt = t
		}
	}
	data, err := json.Marshal(approval)
	if err != nil {
		return fmt.Errorf("marshal approval: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO approvals (session_id, run_id, approval_json, requested_at, approved_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (session_id) DO UPDATE SET approval_json = EXCLUDED.approval_json, approved_at = EXCLUDED.approved_at`,
		approval.SessionID, session.RunID, data, requestedAt, approvedAt)
	return err
}

func dataJSON(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}
