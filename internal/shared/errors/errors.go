package errors

import "fmt"

type Code string

const (
	CodeValidation Code = "VALIDATION"
	CodeContract   Code = "CONTRACT"
	CodeInternal   Code = "INTERNAL"
)

// Error is a transport-neutral error envelope shared by contracts and tests.
type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e Error) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s (%s): %s", e.Code, e.Field, e.Message)
}

func NewValidation(field, message string) Error {
	return Error{Code: CodeValidation, Field: field, Message: message}
}
