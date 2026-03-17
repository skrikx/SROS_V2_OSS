#!/usr/bin/env sh
set -eu

go run ./cmd/sros run --contract examples/runs/minimal_run_contract.json
go run ./cmd/sros pause --reason "runtime smoke pause"
go run ./cmd/sros checkpoint --stage validated
go run ./cmd/sros rollback --reason "runtime smoke rollback"
go run ./cmd/sros status --latest
go run ./cmd/sros run --contract examples/runs/ask_mode_run_contract.json
go run ./cmd/sros status --latest
