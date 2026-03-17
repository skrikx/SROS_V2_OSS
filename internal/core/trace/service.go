package trace

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	ctrace "srosv2/contracts/trace"
	"srosv2/internal/shared/ids"
)

type Service struct {
	Ledger *Ledger
	Writer *Writer
	Reader *Reader
	Query  *Query
	Replay *Replayer
}

func New(root string, now func() time.Time) (*Service, error) {
	ledger, err := NewLedger(root)
	if err != nil {
		return nil, err
	}
	reader := NewReader(ledger)
	return &Service{
		Ledger: ledger,
		Writer: NewWriter(ledger, now),
		Reader: reader,
		Query:  NewQuery(reader),
		Replay: NewReplayer(reader),
	}, nil
}

func (s *Service) Emit(runID ids.RunID, traceID ids.TraceID, spanID ids.SpanID, parent ids.SpanID, kind ctrace.EventType, payload map[string]any) (ctrace.TraceEvent, error) {
	return s.Writer.Event(runID, traceID, spanID, parent, kind, payload)
}

func (s *Service) InspectFromFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read trace input: %w", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("decode trace input: %w", err)
	}
	return payload, nil
}
