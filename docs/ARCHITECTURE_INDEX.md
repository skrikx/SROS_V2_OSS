# Architecture Index

This index maps repository ownership zones only.
Implementation details are intentionally deferred beyond W01.

- `cmd/sros/` - minimal CLI bootstrap entrypoint
- `internal/core/` - core planes ownership boundary
- `internal/fabric/` - fabric/tooling ownership boundary
- `internal/ace/` - ACE orchestration ownership boundary
- `internal/shared/` - common primitives boundary
- `contracts/` - canonical schema/contract boundary
- `tests/` - validation suites by test class
- `artifacts/` - generated runtime/build receipts
- `docs/` - constitutional and workflow documentation
- `examples/` - example usage surface (future workflows)
- `scripts/` - local developer helper scripts (future workflows)
