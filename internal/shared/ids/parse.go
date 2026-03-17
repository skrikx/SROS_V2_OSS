package ids

import (
	"fmt"
	"regexp"
	"strings"
)

var safeIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{2,127}$`)

func IsSafe(value string) bool {
	return safeIDPattern.MatchString(value)
}

func Validate(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fmt.Errorf("id is required")
	}
	if !IsSafe(trimmed) {
		return fmt.Errorf("id %q is not filesystem-safe", value)
	}
	return nil
}

func ParseRunID(value string) (RunID, error) {
	if err := Validate(value); err != nil {
		return "", err
	}
	return RunID(value), nil
}

func ParseTraceID(value string) (TraceID, error) {
	if err := Validate(value); err != nil {
		return "", err
	}
	return TraceID(value), nil
}

func ParseSpanID(value string) (SpanID, error) {
	if err := Validate(value); err != nil {
		return "", err
	}
	return SpanID(value), nil
}

func ParseOperatorID(value string) (OperatorID, error) {
	if err := Validate(value); err != nil {
		return "", err
	}
	return OperatorID(value), nil
}

func ParseCheckpointID(value string) (CheckpointID, error) {
	if err := Validate(value); err != nil {
		return "", err
	}
	return CheckpointID(value), nil
}

func ParseArtifactID(value string) (ArtifactID, error) {
	if err := Validate(value); err != nil {
		return "", err
	}
	return ArtifactID(value), nil
}
