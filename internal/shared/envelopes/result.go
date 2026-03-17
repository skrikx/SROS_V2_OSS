package envelopes

import serrors "srosv2/internal/shared/errors"

// Result is a minimal generic envelope for canonical contract responses.
type Result[T any] struct {
	OK    bool           `json:"ok"`
	Data  *T             `json:"data,omitempty"`
	Error *serrors.Error `json:"error,omitempty"`
	Meta  Meta           `json:"meta"`
}

func Success[T any](data T, meta Meta) Result[T] {
	return Result[T]{
		OK:   true,
		Data: &data,
		Meta: meta,
	}
}

func Failure[T any](err serrors.Error, meta Meta) Result[T] {
	return Result[T]{
		OK:    false,
		Error: &err,
		Meta:  meta,
	}
}
