package orch

import "fmt"

type Queue struct {
	items []WorkUnit
}

func NewQueue() *Queue {
	return &Queue{items: []WorkUnit{}}
}

func (q *Queue) Enqueue(unit WorkUnit) {
	q.items = append(q.items, unit)
}

func (q *Queue) Dequeue() (WorkUnit, error) {
	if len(q.items) == 0 {
		return WorkUnit{}, fmt.Errorf("queue is empty")
	}
	unit := q.items[0]
	q.items = q.items[1:]
	return unit, nil
}

func (q *Queue) Len() int {
	return len(q.items)
}

func (q *Queue) Snapshot() []WorkUnit {
	out := make([]WorkUnit, len(q.items))
	copy(out, q.items)
	return out
}
