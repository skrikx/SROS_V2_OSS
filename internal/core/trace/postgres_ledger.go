package trace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	ctrace "srosv2/contracts/trace"
)

type PostgresLedger struct {
	db *sql.DB
}

func NewPostgresLedger(db *sql.DB) *PostgresLedger {
	return &PostgresLedger{db: db}
}

func (l *PostgresLedger) SaveEvent(ctx context.Context, event ctrace.TraceEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal trace event: %w", err)
	}
	if _, err := l.db.ExecContext(ctx, `
		INSERT INTO trace_events (event_id, run_id, trace_id, span_id, event_type, event_json, occurred_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (event_id) DO UPDATE SET event_json = EXCLUDED.event_json`,
		event.EventID, event.RunID, event.TraceID, event.SpanID, event.EventType, data, event.OccurredAt); err != nil {
		return fmt.Errorf("insert trace event: %w", err)
	}
	if event.ReceiptRef != "" {
		if _, err := l.db.ExecContext(ctx, `INSERT INTO trace_receipt_links (event_id, receipt_id) VALUES ($1,$2)`, event.EventID, event.ReceiptRef); err != nil {
			return fmt.Errorf("insert trace receipt link: %w", err)
		}
	}
	return nil
}

func (l *PostgresLedger) SaveSpan(ctx context.Context, span ctrace.Span) error {
	data, err := json.Marshal(span)
	if err != nil {
		return fmt.Errorf("marshal span: %w", err)
	}
	_, err = l.db.ExecContext(ctx, `
		INSERT INTO evaluation_results (evaluation_id, run_id, result_json, created_at)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (evaluation_id) DO UPDATE SET result_json = EXCLUDED.result_json`,
		span.SpanID, span.TraceID, data, span.StartedAt)
	return err
}
