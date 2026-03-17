# Connector Auth Envelopes

Connector credentials are carried only in auth envelopes.

- Envelope metadata may be inspected.
- Secret material is redacted from CLI output.
- Expired envelopes fail closed.
- Envelope inspection emits trace and receipt hooks through the existing planes.
