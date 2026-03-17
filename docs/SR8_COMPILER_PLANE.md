# SR8 Compiler Plane

SR8 is the compile ingress for SROS V2.

Flow:
1. Parse raw intent input (`--intent` or `--input`).
2. Normalize text into stable compile form.
3. Classify domain and risk.
4. Select compile topology and route class.
5. Validate compile request and run contract constraints.
6. Assemble canonical run contract.
7. Emit compile artifacts (run contract JSON, optional SRXML, compile receipt).

SR8 does not admit or execute runs.
SR9 remains the runtime gate.
