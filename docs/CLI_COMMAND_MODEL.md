# CLI Command Model

SROS V2 uses `cmd/sros` as the canonical local-only operator ingress.

## Families

- Bootstrap/environment: `init`, `bootstrap`, `doctor`, `seed`, `config`
- Compile/run: `compile`, `run`, `plan`, `resume`, `pause`, `checkpoint`, `rollback`
- Inspect/witness: `trace`, `receipts`, `memory`, `mirror`, `inspect`, `status`
- Fabric/capabilities: `tools {list,show,validate,register}`, `connectors list`, `mcp ingest`

## Output Modes

- `text` default for human operators.
- `json` for automation (`--format json` or `SROS_OUTPUT_FORMAT=json`).

## Deferred Semantics

Runtime control commands (`run`, `resume`, `pause`, `checkpoint`, `rollback`, `status`) are wired in W05 through SR9 and the runtime manager.
Other planes that are still pending (trace, memory, mirror, fabric capabilities) fail intentionally with explicit deferred messages and non-zero exit codes.
No command silently embeds later-plane behavior.
