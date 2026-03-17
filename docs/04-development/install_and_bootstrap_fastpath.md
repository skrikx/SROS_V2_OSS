# Install And Bootstrap Fastpath

Fast path:

1. `go build ./cmd/sros`
2. `go run ./cmd/sros verify`
3. `scripts/first_run_smoke.sh`

Optional local Postgres path:

1. `compose.yaml`
2. export `SROS_DATABASE_URL`
3. `scripts/migrate_local.sh`
4. `scripts/seed_local.sh`
