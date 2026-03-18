#!/usr/bin/env bash
set -euo pipefail

STRICT="${SROS_STRICT:-0}"
if [ "${CI:-}" = "true" ] || [ "${CI:-}" = "1" ]; then
  STRICT=1
fi

DB_URL="${SROS_DATABASE_URL:-${DATABASE_URL:-}}"

if [ -z "$DB_URL" ]; then
  if [ "$STRICT" = "1" ]; then
    echo "blocked_by_missing_postgres_runtime: no database url configured" >&2
    exit 1
  fi
  echo "blocked_by_missing_postgres_runtime: no database url configured"
  exit 0
fi

if ! command -v psql >/dev/null 2>&1; then
  if [ "$STRICT" = "1" ]; then
    echo "blocked_by_missing_postgres_runtime: psql not found" >&2
    exit 1
  fi
  echo "blocked_by_missing_postgres_runtime: psql not found"
  exit 0
fi

for file in migrations/*.sql; do
  psql "$DB_URL" -f "$file"
done

echo "migrations_applied"
