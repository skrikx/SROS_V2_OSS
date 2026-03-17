# Run State Machine

Canonical SR9 runtime states:

1. `planned`
2. `approved`
3. `running`
4. `waiting_for_input`
5. `paused`
6. `checkpointed`
7. `failed_safe`
8. `rolled_back`
9. `completed`

Terminal states:

1. `failed_safe`
2. `rolled_back`
3. `completed`

Transition law:

1. Transitions are explicit and deterministic.
2. Illegal transitions return an error.
3. Every state change records reason, timestamp, and session linkage.

Local Ask-mode adaptation:

1. Admission `ask` transitions the session to `waiting_for_input`.
2. SR9 writes a local operator approval artifact.
3. Resume from `waiting_for_input` requires explicit approval input.
