package trace

type EventType string

const (
	EventTypeRunStarted     EventType = "run.started"
	EventTypeRunCompleted   EventType = "run.completed"
	EventTypePolicyDecision EventType = "policy.decision"
	EventTypeMemoryMutation EventType = "memory.mutation"
	EventTypeReceiptLinked  EventType = "receipt.linked"
)
