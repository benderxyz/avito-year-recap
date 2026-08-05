package aggregation

import (
	"math"
	"testing"
)

func TestSanitizeFloatShouldReturnZeroWhenValueIsNaN(t *testing.T) {
	if got := sanitizeFloat(math.NaN()); got != 0 {
		t.Fatalf("expected 0 for NaN, got %v", got)
	}
}

func TestSanitizeFloatShouldReturnZeroWhenValueIsInf(t *testing.T) {
	if got := sanitizeFloat(math.Inf(1)); got != 0 {
		t.Fatalf("expected 0 for +Inf, got %v", got)
	}
}

func TestSanitizeFloatShouldKeepFiniteValueWhenValid(t *testing.T) {
	if got := sanitizeFloat(12.5); got != 12.5 {
		t.Fatalf("expected 12.5, got %v", got)
	}
}
