# Tool Manifest Lifecycle

Lifecycle states:

- `draft`
- `validated`
- `admitted`
- `experimental`
- `active`
- `quarantined`
- `deprecated`
- `disabled`

Rules:

- Draft and validated capabilities are not runnable.
- Admission requires policy binding and harness compatibility.
- Quarantined capabilities remain visible but cannot be selected.
- Deprecated capabilities remain historical only.
- Disabled capabilities remain receiptable but non-runnable.
