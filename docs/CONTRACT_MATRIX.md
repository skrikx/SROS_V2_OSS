# Contract Matrix

| Contract Family | Path | Downstream Owner Workflow |
|---|---|---|
| Shared IDs and primitives | `internal/shared/*` | W03-W10 (cross-cutting) |
| Run contract and checkpoint references | `contracts/runcontract/*` | W04, W05 |
| Trace events and spans | `contracts/trace/*` | W08 |
| Policy decisions and bundles | `contracts/policy/*` | W06 |
| Memory mutation lineage | `contracts/memory/*` | W07 |
| Evidence and receipt envelopes | `contracts/evidence/*` | W08 |
| Checkpoint/release/rollback records | `contracts/release/*` | W10 |
| Canonical JSON schemas | `contracts/jsonschema/*` | W03-W10 |
| Canonical SRXML samples | `contracts/srxml/*` | W04 |
| Proto transport reservation | `contracts/proto/*` | W09-W10 |

W02 owns type law only. Behavior and orchestration are deferred by workflow design.
