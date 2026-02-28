package errors

import (
	"errors"
	"testing"
)

func TestCodeFromError(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		err := NewNotFound("ORDER_NOT_FOUND", "not found")
		code, ok := CodeFromError(err)
		if !ok || code != "ORDER_NOT_FOUND" {
			t.Fatalf("CodeFromError(%v) = %q, %v; want ORDER_NOT_FOUND, true", err, code, ok)
		}
	})
	t.Run("Validation", func(t *testing.T) {
		err := NewValidation("VALIDATION_ERROR", "invalid", []Detail{{Field: "x", Code: "required", Message: "required"}})
		code, ok := CodeFromError(err)
		if !ok || code != "VALIDATION_ERROR" {
			t.Fatalf("CodeFromError(%v) = %q, %v; want VALIDATION_ERROR, true", err, code, ok)
		}
	})
	t.Run("plain_error", func(t *testing.T) {
		err := errors.New("plain")
		code, ok := CodeFromError(err)
		if ok || code != "" {
			t.Fatalf("CodeFromError(plain) = %q, %v; want \"\", false", code, ok)
		}
	})
}
