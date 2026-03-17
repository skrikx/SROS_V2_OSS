package validation

import (
	"strings"
	"time"

	serrors "srosv2/internal/shared/errors"
)

func RequiredString(field, value string) error {
	if strings.TrimSpace(value) == "" {
		err := serrors.NewValidation(field, "is required")
		return err
	}
	return nil
}

func RequiredTime(field string, value time.Time) error {
	if value.IsZero() {
		err := serrors.NewValidation(field, "is required")
		return err
	}
	return nil
}

func RequiredSlice[T any](field string, value []T) error {
	if len(value) == 0 {
		err := serrors.NewValidation(field, "must contain at least one item")
		return err
	}
	return nil
}
