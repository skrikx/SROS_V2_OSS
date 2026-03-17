# Bootstrap (Day 1)

W01 provides the minimum legal local developer flow:

1. Ensure Go is installed.
2. Build bootstrap CLI:
   - `go build ./cmd/sros`
3. Run test baseline:
   - `go test ./...`
4. Run the bootstrap binary:
   - `go run ./cmd/sros`

## Notes

- Current bootstrap is intentionally minimal and compile-safe.
- Command routing and subsystem execution are intentionally deferred to later workflows.
- This flow validates scaffold integrity, not runtime feature completeness.
