#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros examples run
go run ./cmd/sros tools validate --manifest examples/tools/local_patch_manifest.json
echo "seed_complete"
