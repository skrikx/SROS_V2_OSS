package runtime

import (
	"fmt"
	"sort"
)

var legalTransitions = map[SessionState]map[SessionState]bool{
	SessionStatePlanned: {
		SessionStateApproved:        true,
		SessionStateWaitingForInput: true,
		SessionStateFailedSafe:      true,
	},
	SessionStateApproved: {
		SessionStateRunning:         true,
		SessionStateWaitingForInput: true,
		SessionStateFailedSafe:      true,
	},
	SessionStateRunning: {
		SessionStateWaitingForInput: true,
		SessionStatePaused:          true,
		SessionStateCheckpointed:    true,
		SessionStateCompleted:       true,
		SessionStateFailedSafe:      true,
	},
	SessionStateWaitingForInput: {
		SessionStateApproved:   true,
		SessionStateRunning:    true,
		SessionStatePaused:     true,
		SessionStateFailedSafe: true,
	},
	SessionStatePaused: {
		SessionStateRunning:      true,
		SessionStateCheckpointed: true,
		SessionStateRolledBack:   true,
		SessionStateFailedSafe:   true,
	},
	SessionStateCheckpointed: {
		SessionStateRunning:    true,
		SessionStatePaused:     true,
		SessionStateRolledBack: true,
		SessionStateFailedSafe: true,
	},
	SessionStateFailedSafe: {},
	SessionStateRolledBack: {},
	SessionStateCompleted:  {},
}

func CanTransition(from, to SessionState) bool {
	if from == to {
		return true
	}
	next, ok := legalTransitions[from]
	if !ok {
		return false
	}
	return next[to]
}

func EnsureTransition(from, to SessionState) error {
	if CanTransition(from, to) {
		return nil
	}
	allowed := AllowedTransitions(from)
	return fmt.Errorf("illegal state transition %s -> %s (allowed: %v)", from, to, allowed)
}

func AllowedTransitions(from SessionState) []SessionState {
	next := legalTransitions[from]
	out := make([]SessionState, 0, len(next))
	for state := range next {
		out = append(out, state)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func IsTerminal(state SessionState) bool {
	return state == SessionStateFailedSafe || state == SessionStateRolledBack || state == SessionStateCompleted
}
