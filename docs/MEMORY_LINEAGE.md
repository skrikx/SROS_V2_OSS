# Memory Lineage

W07 memory lineage is explicit and trace-compatible.

- Each durable memory mutation includes actor, scope, branch context, timestamp, and lineage reference.
- Parent mutation references are recorded when available.
- This remains compatible with later W08 trace binding without pretending W08 already exists.
