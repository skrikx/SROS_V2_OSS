package trace

import ctrace "srosv2/contracts/trace"

const (
	EventRunStarted      = ctrace.EventTypeRunStarted
	EventRunCompleted    = ctrace.EventTypeRunCompleted
	EventStateTransition = ctrace.EventTypeStateTransition
	EventPolicyDecision  = ctrace.EventTypePolicyDecision
	EventWorkUnit        = ctrace.EventTypeWorkUnit
	EventMemoryMutation  = ctrace.EventTypeMemoryMutation
	EventMirrorWitness   = ctrace.EventTypeMirrorWitness
	EventArtifactLinked  = ctrace.EventTypeArtifactLinked
	EventReceiptLinked   = ctrace.EventTypeReceiptLinked
	EventClosureProof    = ctrace.EventTypeClosureProof
)
