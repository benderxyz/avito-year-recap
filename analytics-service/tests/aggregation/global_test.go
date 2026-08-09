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
	userTotals     map[string]float64
	globalTotals   map[string]float64
	userPresent    map[string]bool
	perUserTotals  map[string][]float64
	float64sCalls  int
	globalSumCalls int
}

func (q *enrichmentQuerier) QueryFloat64(_ context.Context, query string, args ...any) (float64, bool, error) {
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
		q.globalSumCalls++
		for _, arg := range args {
			eventType, ok := arg.(string)
			if !ok {
				continue
			}
			total, exists := q.globalTotals[eventType]
			if !exists {
				continue
			}
			return total, true, nil
		}
		return 0, true, nil
	}

	return 0, false, nil
}

func (q *enrichmentQuerier) QueryFloat64s(_ context.Context, _ string, args ...any) ([]float64, error) {
	q.float64sCalls++
	for _, arg := range args {
		eventType, ok := arg.(string)
		if !ok {
			continue
		}
		if totals, ok := q.perUserTotals[eventType]; ok {
			return append([]float64(nil), totals...), nil
		}
	}
	return lowerPercentileTotals(88, 12, 0, 100), nil
}

func TestServiceShouldEnrichCounterMetricsWithShareAndPercentile(t *testing.T) {
	querier := &enrichmentQuerier{
		userTotals: map[string]float64{
			"item_published": 47,
		},
		globalTotals: map[string]float64{
			"item_published": 1000,
		},
		perUserTotals: map[string][]float64{
			"item_published": lowerPercentileTotals(88, 12, 0, 100),
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to, "UTC")

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
	service := aggregation.NewService(events.NewRegistry(), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to, "UTC")

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
		perUserTotals: map[string][]float64{
			"item_published": {0, 0, 10, 10},
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to, "UTC")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metric := snapshot.Metrics["listingsPublished"]
	if metric.Share != nil {
		t.Fatalf("expected null share when global total is zero, got %v", metric.Share)
	}
	if metric.Percentile == nil || *metric.Percentile != 0 {
		t.Fatalf("expected percentile 0, got %v", metric.Percentile)
	}
}

func TestServiceShouldEnrichMilestoneWhenValueIsPresent(t *testing.T) {
	querier := &enrichmentQuerier{
		userTotals: map[string]float64{
			"first_item_published": 1_700_000_000,
		},
		perUserTotals: map[string][]float64{
			"first_item_published": higherPercentileTotals(88, 12, 1_800_000_000, 1_600_000_000),
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to, "UTC")

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
	service := aggregation.NewService(events.NewRegistry(), &enrichmentQuerier{})
	from := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := service.Metrics(context.Background(), 42, from, to, "UTC")

	if !apperr.IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestServiceShouldReuseCachedGlobalStatsWhenMetricsAreRequestedTwice(t *testing.T) {
	querier := &enrichmentQuerier{
		userTotals: map[string]float64{
			"item_published": 47,
		},
		globalTotals: map[string]float64{
			"item_published": 1000,
		},
		perUserTotals: map[string][]float64{
			"item_published": lowerPercentileTotals(88, 12, 0, 100),
		},
	}
	service := aggregation.NewService(events.NewRegistry(), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	if _, err := service.Metrics(context.Background(), 42, from, to, "UTC"); err != nil {
		t.Fatalf("unexpected first error: %v", err)
	}
	firstGlobalCalls := querier.globalSumCalls
	firstTotalsCalls := querier.float64sCalls

	if _, err := service.Metrics(context.Background(), 99, from, to, "UTC"); err != nil {
		t.Fatalf("unexpected second error: %v", err)
	}

	if querier.globalSumCalls != firstGlobalCalls {
		t.Fatalf("expected cached global totals, got calls %d after %d", querier.globalSumCalls, firstGlobalCalls)
	}
	if querier.float64sCalls != firstTotalsCalls {
		t.Fatalf("expected cached per-user totals, got calls %d after %d", querier.float64sCalls, firstTotalsCalls)
	}
}

func lowerPercentileTotals(belowCount, restCount int, belowValue, restValue float64) []float64 {
	totals := make([]float64, 0, belowCount+restCount)
	for i := 0; i < belowCount; i++ {
		totals = append(totals, belowValue)
	}
	for i := 0; i < restCount; i++ {
		totals = append(totals, restValue)
	}
	return totals
}

func higherPercentileTotals(aboveCount, restCount int, aboveValue, restValue float64) []float64 {
	totals := make([]float64, 0, aboveCount+restCount)
	for i := 0; i < aboveCount; i++ {
		totals = append(totals, aboveValue)
	}
	for i := 0; i < restCount; i++ {
		totals = append(totals, restValue)
	}
	return totals
}
