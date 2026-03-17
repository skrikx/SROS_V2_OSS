#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros verify
go run ./cmd/sros examples run
go run ./cmd/sros trace inspect --input examples/trace/run_trace_min.json
go run ./cmd/sros receipts export --input examples/provenance/receipt_bundle_min.json
go run ./cmd/sros tools search --class tool.local
