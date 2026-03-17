package shared_test

import (
	"testing"
	"time"

	"srosv2/internal/shared/validation"
)

func TestRequiredString(t *testing.T) {
	if err := validation.RequiredString("field", ""); err == nil {
		t.Fatal("expected required string error")
	}
}

func TestEnum(t *testing.T) {
	if err := validation.Enum("mode", "summary", []string{"none", "summary", "full"}); err != nil {
		t.Fatalf("unexpected enum error: %v", err)
	}
	if err := validation.Enum("mode", "bad", []string{"none", "summary", "full"}); err == nil {
		t.Fatal("expected enum error")
	}
}

func TestRequiredTime(t *testing.T) {
	if err := validation.RequiredTime("at", time.Time{}); err == nil {
		t.Fatal("expected required time error")
	}
}
