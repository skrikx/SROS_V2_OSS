# Migrations

Migration files live in `migrations/` and are ordered:

- `0001_initial_schema.sql`
- `0002_indexes_and_constraints.sql`
- `0003_release_and_evidence_tables.sql`

`scripts/migrate_local.sh` applies them with `psql` when a reachable Postgres runtime is available.
