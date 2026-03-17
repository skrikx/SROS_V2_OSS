#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros verify
go run ./cmd/sros tools search --class tool.local
go run ./cmd/sros connectors envelope inspect --input examples/connectors/local_secret_envelope.json
go run ./cmd/sros mcp ingest --input examples/mcp/ingested_remote_capability.json
go run ./cmd/sros release pack --checkpoint cp_smoke
