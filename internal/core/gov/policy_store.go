package gov

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"srosv2/contracts/policy"
)

type PolicyStore struct {
	db *sql.DB
}

func NewPolicyStore(db *sql.DB) *PolicyStore {
	return &PolicyStore{db: db}
}

func (s *PolicyStore) SaveBundle(ctx context.Context, bundle policy.Bundle) error {
	data, err := json.Marshal(bundle)
	if err != nil {
		return fmt.Errorf("marshal policy bundle: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO policy_bundles (bundle_id, bundle_json, loaded_at)
		VALUES ($1,$2,$3)
		ON CONFLICT (bundle_id) DO UPDATE SET bundle_json = EXCLUDED.bundle_json, loaded_at = EXCLUDED.loaded_at`,
		bundle.BundleID, data, time.Now().UTC())
	return err
}

func (s *PolicyStore) SaveDecision(ctx context.Context, decision policy.PolicyDecision) error {
	data, err := json.Marshal(decision)
	if err != nil {
		return fmt.Errorf("marshal policy decision: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO policy_decisions (decision_id, run_id, trace_id, capability, verdict, decision_json, decided_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (decision_id) DO UPDATE SET decision_json = EXCLUDED.decision_json`,
		decision.DecisionID, decision.RunID, decision.TraceID, decision.Capability, decision.Verdict, data, decision.DecidedAt)
	return err
}
