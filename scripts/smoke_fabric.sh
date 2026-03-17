#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros tools --help
go run ./cmd/sros connectors --help
go run ./cmd/sros mcp --help
go run ./cmd/sros tools validate --manifest examples/tools/local_patch_manifest.json
go run ./cmd/sros tools search --class tool.local
go run ./cmd/sros connectors envelope inspect --input examples/connectors/local_secret_envelope.json
go run ./cmd/sros mcp ingest --input examples/mcp/ingested_remote_capability.json
