CREATE TABLE IF NOT EXISTS run_contracts (
  run_id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  workspace_id TEXT NOT NULL,
  contract_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS runs (
  session_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  tenant_id TEXT NOT NULL,
  workspace_id TEXT NOT NULL,
  state TEXT NOT NULL,
  session_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS run_state_transitions (
  id BIGSERIAL PRIMARY KEY,
  session_id TEXT NOT NULL,
  run_id TEXT NOT NULL,
  from_state TEXT NOT NULL,
  to_state TEXT NOT NULL,
  reason TEXT NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS checkpoints (
  checkpoint_id TEXT PRIMARY KEY,
  session_id TEXT NOT NULL,
  run_id TEXT NOT NULL,
  stage TEXT NOT NULL,
  checkpoint_json JSONB NOT NULL,
  recorded_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS approvals (
  session_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  approval_json JSONB NOT NULL,
  requested_at TIMESTAMPTZ NOT NULL,
  approved_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS policy_bundles (
  bundle_id TEXT PRIMARY KEY,
  bundle_json JSONB NOT NULL,
  loaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS policy_decisions (
  decision_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  trace_id TEXT NOT NULL,
  capability TEXT NOT NULL,
  verdict TEXT NOT NULL,
  decision_json JSONB NOT NULL,
  decided_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS memory_nodes (
  record_key TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  workspace_id TEXT NOT NULL,
  branch_id TEXT NOT NULL,
  record_json JSONB NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS memory_mutations (
  mutation_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  session_id TEXT NOT NULL,
  tenant_id TEXT NOT NULL,
  workspace_id TEXT NOT NULL,
  branch_id TEXT NOT NULL,
  kind TEXT NOT NULL,
  mutation_json JSONB NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS session_branches (
  branch_id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  workspace_id TEXT NOT NULL,
  branch_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS session_trees (
  session_id TEXT PRIMARY KEY,
  tree_json JSONB NOT NULL,
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trace_events (
  event_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  trace_id TEXT NOT NULL,
  span_id TEXT NOT NULL,
  event_type TEXT NOT NULL,
  event_json JSONB NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS trace_receipt_links (
  id BIGSERIAL PRIMARY KEY,
  event_id TEXT NOT NULL,
  receipt_id TEXT NOT NULL,
  linked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS receipts (
  receipt_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  kind TEXT NOT NULL,
  receipt_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS evidence_bundles (
  bundle_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  bundle_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS evaluation_results (
  evaluation_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  result_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
