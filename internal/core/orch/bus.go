package orch

import "time"

type Event struct {
	Type      string    `json:"type"`
	SessionID string    `json:"session_id"`
	WorkUnit  string    `json:"work_unit,omitempty"`
	Message   string    `json:"message"`
	At        time.Time `json:"at"`
}

type Bus struct {
	events []Event
}

func NewBus() *Bus {
	return &Bus{events: []Event{}}
}

func (b *Bus) Publish(event Event) {
	b.events = append(b.events, event)
}

func (b *Bus) Events() []Event {
	out := make([]Event, len(b.events))
	copy(out, b.events)
	return out
}
