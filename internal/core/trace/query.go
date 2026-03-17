package trace

import (
	ctrace "srosv2/contracts/trace"
	"srosv2/internal/shared/ids"
)

type Query struct {
	reader *Reader
}

func NewQuery(reader *Reader) *Query {
	return &Query{reader: reader}
}

func (q *Query) ByType(runID ids.RunID, kind ctrace.EventType) ([]ctrace.TraceEvent, error) {
	events, err := q.reader.Events(runID)
	if err != nil {
		return nil, err
	}
	filtered := []ctrace.TraceEvent{}
	for _, event := range events {
		if event.EventType == kind {
			filtered = append(filtered, event)
		}
	}
	return filtered, nil
}
