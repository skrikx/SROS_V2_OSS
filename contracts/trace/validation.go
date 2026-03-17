package trace

import (
	"fmt"

	"srosv2/internal/shared/validation"
)

func ValidateEvent(event TraceEvent) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", event.ContractVersion))
	appendErr(validation.RequiredString("event_id", string(event.EventID)))
	appendErr(validation.RequiredString("trace_id", string(event.TraceID)))
	appendErr(validation.RequiredString("span_id", string(event.SpanID)))
	appendErr(validation.RequiredString("run_id", string(event.RunID)))
	appendErr(validation.Enum("event_type", string(event.EventType), []string{
		"run.started", "run.completed", "policy.decision", "memory.mutation", "receipt.linked",
	}))
	appendErr(validation.RequiredTime("occurred_at", event.OccurredAt))
	return errs
}

func ValidateSpan(span Span) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("trace_id", string(span.TraceID)))
	appendErr(validation.RequiredString("span_id", string(span.SpanID)))
	appendErr(validation.RequiredString("name", span.Name))
	appendErr(validation.RequiredTime("started_at", span.StartedAt))
	if !span.EndedAt.IsZero() && span.EndedAt.Before(span.StartedAt) {
		appendErr(fmt.Errorf("ended_at must be after started_at"))
	}
	return errs
}
