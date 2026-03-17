package mem

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	cmemory "srosv2/contracts/memory"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) SaveRecord(ctx context.Context, record MemoryRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal memory record: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO memory_nodes (record_key, tenant_id, workspace_id, branch_id, record_json, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (record_key) DO UPDATE SET record_json = EXCLUDED.record_json, updated_at = EXCLUDED.updated_at`,
		record.Key, record.Scope.TenantID, record.Scope.WorkspaceID, record.BranchID, data, record.UpdatedAt)
	return err
}

func (s *PostgresStore) SaveMutation(ctx context.Context, mutation cmemory.MemoryMutation) error {
	data, err := json.Marshal(mutation)
	if err != nil {
		return fmt.Errorf("marshal memory mutation: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO memory_mutations (mutation_id, run_id, session_id, tenant_id, workspace_id, branch_id, kind, mutation_json, occurred_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (mutation_id) DO UPDATE SET mutation_json = EXCLUDED.mutation_json`,
		mutation.MutationID, mutation.RunID, mutation.SessionID, mutation.TenantID, mutation.WorkspaceID, mutation.Branch.BranchID, mutation.Kind, data, mutation.OccurredAt)
	return err
}

func (s *PostgresStore) SaveBranch(ctx context.Context, branch cmemory.BranchRecord) error {
	data, err := json.Marshal(branch)
	if err != nil {
		return fmt.Errorf("marshal branch record: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO session_branches (branch_id, tenant_id, workspace_id, branch_json, created_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (branch_id) DO UPDATE SET branch_json = EXCLUDED.branch_json`,
		branch.BranchID, branch.TenantID, branch.WorkspaceID, data, branch.CreatedAt)
	return err
}

func (s *PostgresStore) SaveSessionTree(ctx context.Context, sessionID string, tree []SessionNode) error {
	data, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("marshal session tree: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO session_trees (session_id, tree_json, recorded_at)
		VALUES ($1,$2,$3)
		ON CONFLICT (session_id) DO UPDATE SET tree_json = EXCLUDED.tree_json, recorded_at = EXCLUDED.recorded_at`,
		sessionID, data, time.Now().UTC())
	return err
}
