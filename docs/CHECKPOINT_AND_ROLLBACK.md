# Checkpoint And Rollback

W05 makes checkpoint and rollback runtime-real in local mode.

Checkpoint behavior:

1. `sros checkpoint` creates a contract-aligned checkpoint record.
2. The runtime session transitions to `checkpointed`.
3. The latest checkpoint reference is stored in runtime state.

Rollback behavior:

1. `sros rollback` requires a checkpoint target (explicit or latest).
2. SR9 validates checkpoint existence in the local runtime store.
3. A rollback record is created and the session transitions to `rolled_back`.

Persistence note:

1. W05 uses a minimal local runtime store under `artifacts/runtime`.
2. This store is intentionally replaceable by W10 persistence work.
3. No migration system or final database coupling is introduced in W05.
