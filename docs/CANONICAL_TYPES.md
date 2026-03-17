# Canonical Types

## Principles

- Contracts are transport-neutral.
- IDs are typed where identity matters.
- Validation is explicit for required fields and enum domains.
- Packages in `contracts/*` and `internal/shared/*` do not implement runtime engines.

## IDs

`internal/shared/ids` defines typed aliases for run, trace, span, policy, memory, evidence, and release records.
IDs are validated for filesystem-safe shape using a strict character set.

## Envelopes and Errors

`internal/shared/envelopes` provides neutral metadata and generic result envelopes.
`internal/shared/errors` provides a transport-neutral error shape suitable for CLI and file artifacts.

## Validation Posture

`internal/shared/validation` contains side-effect-free helpers for required values and enum checks.
Contract packages use these helpers to reduce duplicated validation code while preserving ownership boundaries.

## Semantic Alignment

V2 keeps V3 semantic law for contracts while narrowing deployment assumptions to local-only CLI mode.
No hosted control-plane assumptions are present in W02 contract types.
