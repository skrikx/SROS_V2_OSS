# MCP Ingress Posture

MCP ingress in V2 is governed intake only.

- Raw MCP capability metadata is normalized into the SROS manifest shape.
- Normalized manifests are validated and admitted through the registry.
- No direct raw MCP invocation bypasses registry, GOV, or harness posture.
- MCP egress hosting remains deferred.
