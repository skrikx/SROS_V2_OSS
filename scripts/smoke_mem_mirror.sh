#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros memory --help
go run ./cmd/sros mirror --help
go run ./cmd/sros memory recall --input examples/memory/workspace_seed.json
go run ./cmd/sros memory branch --input examples/memory/branch_lineage.json
go run ./cmd/sros mirror witness --input examples/mirror/witness_case.json
go run ./cmd/sros mirror status --input examples/mirror/runtime_snapshot.json
