# SROS V2

[![CI](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/ci.yml/badge.svg)](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/ci.yml)
[![Integration Postgres](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/integration-postgres.yml/badge.svg)](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/integration-postgres.yml)
[![Release Verify](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/release-verify.yml/badge.svg)](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/release-verify.yml)
[![Docs Links](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/docs-links.yml/badge.svg)](https://github.com/skrikx/SROS_V2_OSS/actions/workflows/docs-links.yml)

SROS V2 is a local-only, CLI-first sovereign kernel for governed runs, traces, receipts, memory, replay, and tool-fabric workflows.
Public automation validates build and tests, Postgres integration paths, replay and regression coverage, release-pack verification, and docs link integrity.

It is not a toy shell.
It is the narrowed local profile of SROS V3, with the real chain intact:

`Intent -> SR8 -> run contract / SRXML -> SR9 -> ORCH / GOV / MEM / MIRROR -> TRACE -> PROVENANCE`

## Why try it

In the first run you can get to real proof, not a mock demo:

- a verified front-door readiness report
- a real trace inspection
- a real receipt export
- a governed fabric search
- a local release-pack artifact

## Fast path

```bash
go build ./cmd/sros
go run ./cmd/sros verify
go run ./cmd/sros examples run
go run ./cmd/sros trace inspect --input examples/trace/run_trace_min.json
go run ./cmd/sros receipts export --input examples/provenance/receipt_bundle_min.json
```

Scripted version:

```bash
./scripts/first_run_smoke.sh
```

## Strongest repo proofs

- `examples/showcase/minimal_governed_run/`
- `examples/showcase/trace_and_receipt_walkthrough/`
- `examples/showcase/memory_branch_and_rewind/`
- `examples/showcase/tool_registry_search_and_invoke/`
- `examples/showcase/replay_and_closure_proof/`
- `artifacts/showcase/`

## Operator commands

- `sros verify`
- `sros examples catalog`
- `sros inspect`
- `sros trace inspect`
- `sros receipts export`
- `sros test first-run`
- `sros release pack`

## Storage posture

- PostgreSQL for transactional and append-style local persistence when configured
- `artifacts/` for emitted receipts, bundles, releases, replays, and showcase packs

## Start reading

- `docs/01-intro/first_run.md`
- `docs/01-intro/time_to_value.md`
- `docs/03-cli/ux_and_outputs.md`
- `docs/09-examples/killer_examples.md`
- `CONTRIBUTING.md`
