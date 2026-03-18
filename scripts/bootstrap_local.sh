#!/usr/bin/env bash
set -euo pipefail

STRICT="${SROS_STRICT:-0}"
if [ "${CI:-}" = "true" ] || [ "${CI:-}" = "1" ]; then
  STRICT=1
fi

mkdir -p artifacts/receipts artifacts/bundles artifacts/releases artifacts/replays

if command -v docker >/dev/null 2>&1 && [ -f compose.yaml ]; then
  if ! docker compose up -d postgres >/dev/null 2>&1; then
    if [ "$STRICT" = "1" ]; then
      echo "bootstrap_failed: docker compose could not start postgres" >&2
      exit 1
    fi
    echo "bootstrap_warning: docker compose could not start postgres (continuing in local-friendly mode)"
  fi
fi

go run ./cmd/sros scaffold
go run ./cmd/sros verify
