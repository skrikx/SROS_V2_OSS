CREATE INDEX IF NOT EXISTS idx_runs_run_id ON runs(run_id);
CREATE INDEX IF NOT EXISTS idx_run_state_transitions_session_id ON run_state_transitions(session_id);
CREATE INDEX IF NOT EXISTS idx_checkpoints_run_id ON checkpoints(run_id);
CREATE INDEX IF NOT EXISTS idx_policy_decisions_run_id ON policy_decisions(run_id);
CREATE INDEX IF NOT EXISTS idx_memory_mutations_run_id ON memory_mutations(run_id);
CREATE INDEX IF NOT EXISTS idx_memory_mutations_branch_id ON memory_mutations(branch_id);
CREATE INDEX IF NOT EXISTS idx_trace_events_run_id ON trace_events(run_id);
CREATE INDEX IF NOT EXISTS idx_trace_events_trace_id ON trace_events(trace_id);
CREATE INDEX IF NOT EXISTS idx_receipts_run_id ON receipts(run_id);

ALTER TABLE runs
  ADD CONSTRAINT runs_local_semantics CHECK (tenant_id <> '' AND workspace_id <> '');

ALTER TABLE memory_nodes
  ADD CONSTRAINT memory_nodes_local_semantics CHECK (tenant_id <> '' AND workspace_id <> '');

ALTER TABLE memory_mutations
  ADD CONSTRAINT memory_mutations_local_semantics CHECK (tenant_id <> '' AND workspace_id <> '');
