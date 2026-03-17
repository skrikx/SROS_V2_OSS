# MEM Plane

MEM owns explicit local continuity.

- Workspace memory is persisted as local records plus an explicit mutation ledger.
- Every durable mutation records tenant, workspace, operator, scope, branch, and lineage reference.
- Branch and rewind operate on memory continuity only.
- Pruning and compaction are separate operations with different semantics.
