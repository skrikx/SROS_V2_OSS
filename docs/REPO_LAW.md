# Repo Law (W01)

## Constitutional Constraints

- V2 is the local-only, CLI-first sovereign profile of V3.
- There is one architecture and one canonical chain.
- Scaffold before logic is mandatory.

## Hard Prohibitions

- No `cmd/srosd` daemon path.
- No enterprise shell widening.
- No external HTTP API/UI control plane.
- No subsystem business logic in W01 scaffolding packages.

## Workflow Boundary

W01 owns repository scaffold, bootstrap entrypoint, and constitution binding docs only.
All subsystem behavior is deferred to W02-W10.
