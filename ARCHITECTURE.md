# Architecture

SROS V2 is not a forked architecture. It is the local runtime envelope of SROS V3.

## Canonical chain

Intent -> SR8 -> run contract / SRXML -> SR9 -> ORCH / GOV / MEM / MIRROR -> TRACE -> PROVENANCE

## Persistence posture

- PostgreSQL backs durable runtime, policy, memory, trace, evidence, and release classes
- `artifacts/` stores emitted receipts, bundles, release packs, and replay artifacts

## Ownership zones

- `cmd/sros` - operator CLI
- `internal/core` - kernel planes and persistence seams
- `internal/fabric` - governed capability fabric
- `contracts` - canonical contract truth
- `docs`, `examples`, `scripts`, `tests` - developer and OSS support surface

## Exclusions

- no hosted APIs
- no daemon-first shell
- no enterprise auth, RBAC, or admin UI
