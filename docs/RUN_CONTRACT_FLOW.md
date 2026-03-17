# Run Contract Flow

W04 produces canonical run contracts through SR8.

Boundary handoff:

- SR8: parse, normalize, classify, topology, validate, assemble, receipt.
- SR9 (deferred to W05): admission and runtime state transitions.

SR8 output is sufficient for SR9 admission without re-reading raw intent.
