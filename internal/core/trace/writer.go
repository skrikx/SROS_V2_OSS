package trace

import (
	"fmt"
	"time"

	ctrace "srosv2/contracts/trace"
	"srosv2/internal/shared/ids"
)

type Writer struct {
	ledger *Ledger
	now    func() time.Time
}

func NewWriter(ledger *Ledger, now func() time.Time) *Writer {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Writer{ledger: ledger, now: now}
}

func (w *Writer) BeginSpan(traceID ids.TraceID, parent ids.SpanID, name string, attrs map[string]string) (ctrace.Span, error) {
	span := ctrace.Span{
		TraceID:      traceID,
		SpanID:       ids.SpanID("span_" + shortHash(fmt.Sprintf("%s|%s|%d", traceID, name, w.now().UnixNano()))),
		ParentSpanID: parent,
		Name:         name,
		StartedAt:    w.now().UTC(),
		Attributes:   attrs,
	}
	return span, w.ledger.SaveSpan(span)
}

func (w *Writer) EndSpan(span ctrace.Span) error {
	span.EndedAt = w.now().UTC()
	return w.ledger.SaveSpan(span)
}

func (w *Writer) Event(runID ids.RunID, traceID ids.TraceID, spanID ids.SpanID, parent ids.SpanID, kind ctrace.EventType, payload map[string]any) (ctrace.TraceEvent, error) {
	if spanID == "" {
		spanID = ids.SpanID("span_" + shortHash(fmt.Sprintf("%s|%s|span", runID, kind)))
	}
	event := ctrace.TraceEvent{
		ContractVersion: "v2.0",
		EventID:         ids.EventID("evt_" + shortHash(fmt.Sprintf("%s|%s|%d", runID, kind, w.now().UnixNano()))),
		TraceID:         traceID,
		SpanID:          spanID,
		ParentSpanID:    parent,
		RunID:           runID,
		EventType:       kind,
		OccurredAt:      w.now().UTC(),
		Payload:         payload,
	}
	return event, w.ledger.Append(event)
}
