#!/usr/bin/env bash
set -euo pipefail

mkdir -p artifacts/showcase
cp examples/catalog.json artifacts/showcase/example_catalog.json
cat > artifacts/showcase/first_run_snapshot.json <<'EOF'
{
  "commands": [
    "go run ./cmd/sros verify",
    "go run ./cmd/sros examples run",
    "go run ./cmd/sros trace inspect --input examples/trace/run_trace_min.json",
    "go run ./cmd/sros receipts export --input examples/provenance/receipt_bundle_min.json"
  ],
  "generated_by": "scripts/build_showcase_pack.sh"
}
EOF
cat > artifacts/showcase/share_pack_manifest.json <<'EOF'
{
  "artifacts": [
    "artifacts/showcase/example_catalog.json",
    "artifacts/showcase/first_run_snapshot.json",
    "examples/provenance/receipt_bundle_min.json",
    "examples/trace/run_trace_min.json",
    "examples/releases/release_baseline_min.json"
  ]
}
EOF
echo "showcase_pack_built"
