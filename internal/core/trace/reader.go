package trace

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	ctrace "srosv2/contracts/trace"
	"srosv2/internal/shared/ids"
)

type Reader struct {
	ledger *Ledger
}

func NewReader(ledger *Ledger) *Reader {
	return &Reader{ledger: ledger}
}

func (r *Reader) Events(runID ids.RunID) ([]ctrace.TraceEvent, error) {
	path := filepath.Join(r.ledger.Root(), "events", string(runID)+".jsonl")
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open trace events: %w", err)
	}
	defer file.Close()
	events := []ctrace.TraceEvent{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event ctrace.TraceEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf("decode trace event: %w", err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan trace events: %w", err)
	}
	return events, nil
}

func (r *Reader) Span(spanID ids.SpanID) (ctrace.Span, error) {
	path := filepath.Join(r.ledger.Root(), "spans", string(spanID)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return ctrace.Span{}, fmt.Errorf("read span: %w", err)
	}
	var span ctrace.Span
	if err := json.Unmarshal(data, &span); err != nil {
		return ctrace.Span{}, fmt.Errorf("decode span: %w", err)
	}
	return span, nil
}
