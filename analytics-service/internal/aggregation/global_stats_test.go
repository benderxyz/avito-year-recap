package aggregation

import (
	"context"
	"testing"
	"time"
)

type countingQuerier struct {
	floatCalls  int
	floatsCalls int
	floatValue  float64
	floatsValue []float64
}

func (q *countingQuerier) QueryFloat64(_ context.Context, _ string, _ ...any) (float64, bool, error) {
	q.floatCalls++
	return q.floatValue, true, nil
}

func (q *countingQuerier) QueryFloat64s(_ context.Context, _ string, _ ...any) ([]float64, error) {
	q.floatsCalls++
	return q.floatsValue, nil
}

func TestGlobalStatsShouldReuseCachedCounterTotalWhenRequestedAgain(t *testing.T) {
	querier := &countingQuerier{floatValue: 1000}
	stats := newGlobalStats(querier, time.Minute)
	req := GlobalRequest{
		EventType: "item_published",
		From:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	first, err := stats.counterTotal(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := stats.counterTotal(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first.Value != 1000 || second.Value != 1000 {
		t.Fatalf("expected cached value 1000, got %v and %v", first.Value, second.Value)
	}
	if querier.floatCalls != 1 {
		t.Fatalf("expected one db call, got %d", querier.floatCalls)
	}
}

func TestGlobalStatsShouldReuseCachedPerUserTotalsWhenRequestedAgain(t *testing.T) {
	querier := &countingQuerier{floatsValue: []float64{1, 2, 3}}
	stats := newGlobalStats(querier, time.Minute)
	req := GlobalRequest{
		EventType: "item_published",
		From:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	first, err := stats.counterTotals(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := stats.counterTotals(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(first) != 3 || len(second) != 3 {
		t.Fatalf("expected cached totals length 3, got %d and %d", len(first), len(second))
	}
	if querier.floatsCalls != 1 {
		t.Fatalf("expected one db call, got %d", querier.floatsCalls)
	}
}
