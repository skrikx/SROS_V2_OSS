# Fabric Overview

SROS V2 fabric is the governed capability spine for local tools, connectors, and MCP-ingested capabilities.

- Registry is the discovery and admission truth.
- Manifests carry lifecycle, sandbox, auth, and trust-boundary data.
- Harness profiles keep local-safe execution explicit.
- Connectors require auth envelopes.
- MCP ingress normalizes external capability metadata into manifest shape before admission.
- Capability use stays linked to GOV, TRACE, and PROVENANCE through narrow hooks.
