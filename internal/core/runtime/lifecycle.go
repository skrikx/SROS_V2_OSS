package runtime

import (
	"strings"
	"time"
)

func Transition(session *RuntimeSession, to SessionState, reason string, at time.Time) error {
	if err := EnsureTransition(session.State, to); err != nil {
		return err
	}

	if strings.TrimSpace(reason) == "" {
		reason = "state transition"
	}

	from := session.State
	session.State = to
	session.Reason = reason
	session.UpdatedAt = at.UTC()
	session.History = append(session.History, LifecycleEvent{
		From:   from,
		To:     to,
		Reason: reason,
		At:     at.UTC(),
	})
	return nil
}
