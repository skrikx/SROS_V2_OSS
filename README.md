# SROS V2 (Rebuild)

SROS V2 is the local-only, CLI-first sovereign profile of SROS V3.
This repository is rebuilt to preserve one architecture and one canonical workflow truth:

Intent -> SR8 -> run contract / SRXML -> SR9 -> ORCH / GOV / MEM / MIRROR -> traces -> receipts.

## W01 Scope

This scaffold binds the repository to the v2 constitution and freezes ownership zones for later workflows.
W01 intentionally includes only:

- canonical directory structure
- minimal `cmd/sros` compile-safe bootstrap entry
- constitution-bound docs and repo law

W01 intentionally excludes subsystem logic, command families, daemon/API surfaces, and enterprise shell widening.

## Local-only Law

- No second architecture
- No `cmd/srosd` daemon scaffold
- No hosted/API/web control plane
- No privileged bypass of governance chain

## Quick Start

```bash
go build ./cmd/sros
go test ./...
./sros
```

See:

- `ARCHITECTURE.md`
- `docs/ARCHITECTURE_INDEX.md`
- `docs/REPO_LAW.md`
- `docs/BOOTSTRAP.md`
