# Architecture (W01 Baseline)

## Positioning

SROS V2 is not a legacy fork. It is the local-only, CLI-first profile of SROS V3.
The architecture remains singular across versions; V2 narrows runtime envelope but does not fork topology.

## Canonical Processing Chain

Intent -> SR8 -> run contract / SRXML -> SR9 -> ORCH / GOV / MEM / MIRROR -> traces -> receipts

## Ownership Zones (Scaffold)

- `cmd/sros` - minimal root CLI bootstrap
- `internal/core` - inherited core planes (boot, runtime, orch, gov, mem, mirror, trace, provenance, sr8, sr9)
- `internal/fabric` - local tool fabric surface (registry, harness, connectors, shell, patch, fileio, web, browser, mcpclient)
- `internal/ace` - ACE planes (intake, routing, delegation, research, governance, evaluation, skills, promptunits)
- `internal/shared` - shared primitives/config/error/id/validation envelope
- `contracts` - canonical run/trace/policy/memory/release contract surfaces
- `artifacts`, `tests`, `docs`, `examples`, `scripts` - lifecycle support zones

## Non-Widening Constraint

W01 does not add command families, runtime engines, subsystem logic, daemon surfaces, or enterprise shell width.
