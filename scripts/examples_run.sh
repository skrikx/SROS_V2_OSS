#!/usr/bin/env bash
set -euo pipefail

go run ./cmd/sros examples run
go run ./cmd/sros examples catalog
