package trace

type EventType string

const (
	EventTypeRunStarted      EventType = "run.started"
	EventTypeRunCompleted    EventType = "run.completed"
	EventTypeStateTransition EventType = "state.transition"
	EventTypePolicyDecision  EventType = "policy.decision"
	EventTypeWorkUnit        EventType = "orch.work_unit"
	EventTypeMemoryMutation  EventType = "memory.mutation"
	EventTypeMirrorWitness   EventType = "mirror.witness"
	EventTypeArtifactLinked  EventType = "artifact.linked"
	EventTypeReceiptLinked   EventType = "receipt.linked"
	EventTypeClosureProof    EventType = "closure.proof"
)
