# First Run

The fastest truthful path to value is:

1. `go build ./cmd/sros`
2. `go run ./cmd/sros verify`
3. `go run ./cmd/sros examples run`
4. `go run ./cmd/sros trace inspect --input examples/trace/run_trace_min.json`
5. `go run ./cmd/sros receipts export --input examples/provenance/receipt_bundle_min.json`

If you want the scripted path:

- `scripts/first_run_smoke.sh`

Proof:

- real trace artifact
- real receipt bundle
- real governed fabric search
- real release-pack surface
