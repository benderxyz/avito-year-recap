package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"analytics-service/internal/aggregation"
	"analytics-service/internal/api"
)

func TestParseRangeShouldRejectPartialBoundsWhenOnlyFromIsSet(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/users/1/metrics?from=2026-01-01T00:00:00Z", nil)

	_, _, err := api.ParseRange(req, "UTC")

	if err == nil {
		t.Fatal("expected error when only from is set")
	}
}

func TestParseRangeShouldUseTimezoneYearWhenBoundsAreOmitted(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/users/1/metrics", nil)

	from, to, err := api.ParseRange(req, "Europe/Moscow")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !from.Before(to) {
		t.Fatalf("expected from before to, got %v .. %v", from, to)
	}
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	expectedYear := time.Now().In(loc).Year()
	if from.In(loc).Year() != expectedYear {
		t.Fatalf("expected local year %d, got %d", expectedYear, from.In(loc).Year())
	}
}

func TestYearRangeMatchesAggregationHelper(t *testing.T) {
	now := time.Date(2026, 8, 5, 10, 0, 0, 0, time.UTC)
	from, to, err := aggregation.YearRangeInTimezone(now, "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !from.Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected from %v", from)
	}
	if !to.Equal(now) {
		t.Fatalf("unexpected to %v", to)
	}
}
