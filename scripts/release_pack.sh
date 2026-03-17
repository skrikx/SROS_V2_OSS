#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros release pack --checkpoint cp_release
go run ./cmd/sros verify
