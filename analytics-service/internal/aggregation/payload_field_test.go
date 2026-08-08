package aggregation

import "testing"

func TestSanitizePayloadFieldShouldKeepValidIdentifierWhenFieldMatchesPattern(t *testing.T) {
	if got := sanitizePayloadField("category"); got != "category" {
		t.Fatalf("expected category, got %s", got)
	}
}

func TestSanitizePayloadFieldShouldFallbackToValueWhenFieldIsEmpty(t *testing.T) {
	if got := sanitizePayloadField(""); got != "value" {
		t.Fatalf("expected value, got %s", got)
	}
}

func TestSanitizePayloadFieldShouldFallbackToValueWhenFieldContainsInjection(t *testing.T) {
	if got := sanitizePayloadField("x'); DROP TABLE events; --"); got != "value" {
		t.Fatalf("expected value, got %s", got)
	}
}
