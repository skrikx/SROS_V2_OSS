#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros trace --help
go run ./cmd/sros receipts --help
go run ./cmd/sros trace inspect --input examples/trace/run_trace_min.json
go run ./cmd/sros trace replay --input examples/trace/replay_case.json
go run ./cmd/sros receipts export --input examples/provenance/receipt_bundle_min.json
go run ./cmd/sros receipts closure --input examples/provenance/closure_proof_min.json
