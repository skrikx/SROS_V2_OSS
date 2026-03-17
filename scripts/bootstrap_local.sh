#!/usr/bin/env bash
set -euo pipefail

mkdir -p artifacts/receipts artifacts/bundles artifacts/releases artifacts/replays

if command -v docker >/dev/null 2>&1 && [ -f compose.yaml ]; then
  docker compose up -d postgres >/dev/null 2>&1 || true
fi

go run ./cmd/sros scaffold
go run ./cmd/sros verify
