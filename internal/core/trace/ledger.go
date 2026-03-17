package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	ctrace "srosv2/contracts/trace"
)

type Ledger struct {
	root string
	pg   *PostgresLedger
}

func NewLedger(root string) (*Ledger, error) {
	if root == "" {
		root = filepath.Join("artifacts", "trace")
	}
	for _, rel := range []string{"events", "spans"} {
		if err := os.MkdirAll(filepath.Join(root, rel), 0o755); err != nil {
			return nil, fmt.Errorf("create trace ledger root: %w", err)
		}
	}
	return &Ledger{root: root}, nil
}

func (l *Ledger) Root() string { return l.root }

func (l *Ledger) SetPostgresLedger(pg *PostgresLedger) {
	l.pg = pg
}

func (l *Ledger) Append(event ctrace.TraceEvent) error {
	if errs := ctrace.ValidateEvent(event); len(errs) > 0 {
		return errs[0]
	}
	path := filepath.Join(l.root, "events", string(event.RunID)+".jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open trace ledger: %w", err)
	}
	defer f.Close()
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal trace event: %w", err)
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("append trace event: %w", err)
	}
	if l.pg != nil {
		return l.pg.SaveEvent(context.Background(), event)
	}
	return nil
}

func (l *Ledger) SaveSpan(span ctrace.Span) error {
	data, err := json.MarshalIndent(span, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal span: %w", err)
	}
	path := filepath.Join(l.root, "spans", string(span.SpanID)+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write span: %w", err)
	}
	if l.pg != nil {
		return l.pg.SaveSpan(context.Background(), span)
	}
	return nil
}
