#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros replay --input examples/traces/replay_case_min.json
echo "replay_complete"
