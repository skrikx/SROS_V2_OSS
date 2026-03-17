CREATE TABLE IF NOT EXISTS releases (
  release_id TEXT PRIMARY KEY,
  checkpoint_id TEXT NOT NULL,
  target_stage TEXT NOT NULL,
  release_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS promotion_gate_results (
  gate_id TEXT PRIMARY KEY,
  release_id TEXT NOT NULL,
  status TEXT NOT NULL,
  result_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS rollback_records (
  rollback_id TEXT PRIMARY KEY,
  release_id TEXT NOT NULL,
  target_checkpoint_id TEXT NOT NULL,
  rollback_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS artifact_provenance (
  artifact_id TEXT PRIMARY KEY,
  run_id TEXT NOT NULL,
  source_kind TEXT NOT NULL,
  artifact_json JSONB NOT NULL,
  linked_at TIMESTAMPTZ NOT NULL
);
