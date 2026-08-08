package aggregation_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"analytics-service/internal/aggregation"
	"analytics-service/internal/apperr"
	"analytics-service/internal/events"
)

type enrichmentQuerier struct {
	userTotals   map[string]float64
	globalTotals map[string]float64
	userPresent  map[string]bool
}

func (q *enrichmentQuerier) QueryFloat64(_ context.Context, query string, args ...any) (float64, bool, error) {
	if strings.Contains(query, "countIf(total < ?)") || strings.Contains(query, "countIf(total > ?)") {
		userValue, ok := args[0].(float64)
		if !ok {
			return 0, false, nil
		}
		if userValue >= 40 {
			return 88, true, nil
		}
		return 50, true, nil
	}

	if strings.Contains(query, "user_id = ?") {
		eventType, ok := args[1].(string)
		if !ok {
			return 0, false, nil
		}
		if present, ok := q.userPresent[eventType]; ok && !present {
			return 0, false, nil
		}
		total, exists := q.userTotals[eventType]
		if !exists {
			return 0, true, nil
		}
		return total, true, nil
	}

	if strings.Contains(query, "sum(user_total)") || strings.Contains(query, "sum(value)") {
		eventType, ok := args[0].(string)
		if !ok {
			return 0, false, nil
		}
		total, exists := q.globalTotals[eventType]
		if !exists {
			return 0, true, nil
		}
		return total, true, nil
	}

	return 0, false, nil
}

func TestServiceShouldEnrichCounterMetricsWithShareAndPercentile(t *testing.T) {
	querier := &enrichmentQuerier{
		userTotals: map[string]float64{
			"item_published": 47,
		},
		globalTotals: map[string]float64{
			"item_published": 1000,
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier, staticTimezoneResolver{timezone: "UTC"})
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metric := snapshot.Metrics["listingsPublished"]
	if metric.Value == nil || *metric.Value != 47 {
		t.Fatalf("expected listingsPublished value 47, got %v", metric.Value)
	}
	if metric.Share == nil || *metric.Share != 4.7 {
		t.Fatalf("expected listingsPublished share 4.7, got %v", metric.Share)
	}
	if metric.Percentile == nil || *metric.Percentile != 88 {
		t.Fatalf("expected listingsPublished percentile 88, got %v", metric.Percentile)
	}
}

func TestServiceShouldReturnAllNullFieldsWhenSparseGaugeIsAbsent(t *testing.T) {
	querier := &enrichmentQuerier{
		userPresent: map[string]bool{
			"active_items_count": false,
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier, staticTimezoneResolver{timezone: "UTC"})
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metric := snapshot.Metrics["activeListings"]
	if metric.Value != nil || metric.Percentile != nil || metric.Share != nil {
		t.Fatalf("expected all-null sparse gauge metric, got %+v", metric)
	}
}

func TestServiceShouldReturnNullShareWhenGlobalCounterTotalIsZero(t *testing.T) {
	querier := &enrichmentQuerier{
		userTotals: map[string]float64{
			"item_published": 0,
		},
		globalTotals: map[string]float64{
			"item_published": 0,
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier, staticTimezoneResolver{timezone: "UTC"})
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metric := snapshot.Metrics["listingsPublished"]
	if metric.Share != nil {
		t.Fatalf("expected null share when global total is zero, got %v", metric.Share)
	}
	if metric.Percentile == nil || *metric.Percentile != 50 {
		t.Fatalf("expected percentile 50, got %v", metric.Percentile)
	}
}

func TestServiceShouldEnrichMilestoneWhenValueIsPresent(t *testing.T) {
	querier := &enrichmentQuerier{
		userTotals: map[string]float64{
			"first_item_published": 1_700_000_000,
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier, staticTimezoneResolver{timezone: "UTC"})
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metric := snapshot.Metrics["firstListingAt"]
	if metric.Value == nil || *metric.Value != 1_700_000_000 {
		t.Fatalf("expected firstListingAt value, got %v", metric.Value)
	}
	if metric.Share != nil {
		t.Fatalf("expected null share for milestone, got %v", metric.Share)
	}
	if metric.Percentile == nil || *metric.Percentile != 88 {
		t.Fatalf("expected milestone percentile 88, got %v", metric.Percentile)
	}
}

func TestServiceShouldFailEnrichmentWhenFromIsNotBeforeTo(t *testing.T) {
	service := aggregation.NewService(events.NewRegistry(), &enrichmentQuerier{}, staticTimezoneResolver{timezone: "UTC"})
	from := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := service.Metrics(context.Background(), 42, from, to)

	if !apperr.IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
