package apperr_test

import (
	"errors"
	"testing"

	"analytics-service/internal/apperr"
)

func TestIsValidationShouldDetectValidationErrors(t *testing.T) {
	err := apperr.Validation("bad input")

	if !apperr.IsValidation(err) {
		t.Fatal("expected validation error")
	}
}

func TestIsValidationShouldIgnoreRegularErrors(t *testing.T) {
	err := errors.New("clickhouse down")

	if apperr.IsValidation(err) {
		t.Fatal("did not expect validation error")
	}
}
