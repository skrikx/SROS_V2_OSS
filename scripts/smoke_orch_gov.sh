#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros run --contract examples/runs/governed_runtime_session.json
go run ./cmd/sros checkpoint --latest
go run ./cmd/sros status --latest
go run ./cmd/sros tools validate --policy examples/policy/local_default_policy.json
go run ./cmd/sros tools validate --policy examples/policy/high_risk_patch_policy.json
