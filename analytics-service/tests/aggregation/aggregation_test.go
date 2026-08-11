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

type recordingQuerier struct {
	lastQuery string
	lastArgs  []any
	result    float64
	present   bool
	err       error
}

func (q *recordingQuerier) QueryFloat64(_ context.Context, query string, args ...any) (float64, bool, error) {
	q.lastQuery = query
	q.lastArgs = args
	return q.result, q.present, q.err
}

func (q *recordingQuerier) QueryFloat64s(_ context.Context, query string, args ...any) ([]float64, error) {
	q.lastQuery = query
	q.lastArgs = args
	if q.err != nil {
		return nil, q.err
	}
	if !q.present {
		return nil, nil
	}
	return []float64{q.result}, nil
}

func sampleRequest(eventType string) aggregation.Request {
	return aggregation.Request{
		UserID:    42,
		EventType: eventType,
		From:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
		Timezone:  "Europe/Moscow",
	}
}

func TestFactoryShouldReturnCounterAggregatorWhenCategoryIsCounter(t *testing.T) {
	querier := &recordingQuerier{}

	agg, err := aggregation.NewAggregator(querier, events.CategoryConfig{Category: events.CategoryCounter})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := agg.(*aggregation.CounterAggregator); !ok {
		t.Fatalf("expected *CounterAggregator, got %T", agg)
	}
}

func TestFactoryShouldReturnUniqueAggregatorWhenCategoryIsUnique(t *testing.T) {
	querier := &recordingQuerier{}

	agg, err := aggregation.NewAggregator(querier, events.CategoryConfig{
		Category:    events.CategoryUnique,
		UniqueMode:  events.UniqueModePayload,
		UniqueField: "category",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	unique, ok := agg.(*aggregation.UniqueAggregator)
	if !ok {
		t.Fatalf("expected *UniqueAggregator, got %T", agg)
	}
	if unique.PayloadField() != "category" {
		t.Fatalf("expected payload field category, got %s", unique.PayloadField())
	}
}

func TestFactoryShouldFailWhenCategoryIsUnsupported(t *testing.T) {
	querier := &recordingQuerier{}

	_, err := aggregation.NewAggregator(querier, events.CategoryConfig{Category: events.EventCategory("unknown")})

	if err == nil {
		t.Fatal("expected error for unsupported category")
	}
}

func TestCounterAggregatorShouldUseSumQueryWhenAggregating(t *testing.T) {
	querier := &recordingQuerier{result: 47, present: true}
	agg := aggregation.NewCounterAggregator(querier)

	result, err := agg.Aggregate(context.Background(), sampleRequest("item_published"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Present || result.Value != 47 {
		t.Fatalf("expected present 47, got %+v", result)
	}
	if !strings.Contains(querier.lastQuery, "sum(value)") {
		t.Fatalf("expected sum(value) in query, got %s", querier.lastQuery)
	}
}

func TestGaugeAggregatorShouldUseArgMaxQueryWhenAggregating(t *testing.T) {
	querier := &recordingQuerier{result: 5, present: true}
	agg := aggregation.NewGaugeAggregator(querier)

	result, err := agg.Aggregate(context.Background(), sampleRequest("active_items_count"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Present || result.Value != 5 {
		t.Fatalf("expected present 5, got %+v", result)
	}
	if !strings.Contains(querier.lastQuery, "argMax(value, occurred_at)") {
		t.Fatalf("expected argMax query, got %s", querier.lastQuery)
	}
}

func TestMilestoneAggregatorShouldUseMinOccurredAtWhenAggregating(t *testing.T) {
	querier := &recordingQuerier{result: 1_700_000_000, present: true}
	agg := aggregation.NewMilestoneAggregator(querier)

	result, err := agg.Aggregate(context.Background(), sampleRequest("first_item_published"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Present || result.Value != 1_700_000_000 {
		t.Fatalf("expected present unix timestamp, got %+v", result)
	}
	if !strings.Contains(querier.lastQuery, "min(occurred_at)") {
		t.Fatalf("expected min(occurred_at) query, got %s", querier.lastQuery)
	}
	if !strings.Contains(querier.lastQuery, "count() = 0") {
		t.Fatalf("expected count() absence check, got %s", querier.lastQuery)
	}
}

func TestUniqueAggregatorShouldIgnoreEmptyPayloadFieldWhenAggregating(t *testing.T) {
	querier := &recordingQuerier{result: 2, present: true}
	agg := aggregation.NewUniqueAggregator(querier, events.UniqueModePayload, "category")

	result, err := agg.Aggregate(context.Background(), sampleRequest("category_opened"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Present || result.Value != 2 {
		t.Fatalf("expected present 2, got %+v", result)
	}
	if !strings.Contains(querier.lastQuery, "JSONExtractString(payload, 'category') != ''") {
		t.Fatalf("expected empty payload filter, got %s", querier.lastQuery)
	}
}

func TestUniqueAggregatorShouldFallbackToValueWhenPayloadFieldIsInvalid(t *testing.T) {
	agg := aggregation.NewUniqueAggregator(&recordingQuerier{}, events.UniqueModePayload, "x'); DROP TABLE events; --")

	if agg.PayloadField() != "value" {
		t.Fatalf("expected fallback field value, got %s", agg.PayloadField())
	}
}

func TestUniqueAggregatorShouldCountLocalDaysWhenModeIsDay(t *testing.T) {
	querier := &recordingQuerier{result: 214, present: true}
	agg := aggregation.NewUniqueAggregator(querier, events.UniqueModeDay, "")

	result, err := agg.Aggregate(context.Background(), sampleRequest("day_active"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Present || result.Value != 214 {
		t.Fatalf("expected present 214, got %+v", result)
	}
	if !strings.Contains(querier.lastQuery, "toDate(toTimeZone(occurred_at, ?))") {
		t.Fatalf("expected local day unique expression, got %s", querier.lastQuery)
	}
	if len(querier.lastArgs) == 0 || querier.lastArgs[0] != "Europe/Moscow" {
		t.Fatalf("expected timezone arg Europe/Moscow, got %#v", querier.lastArgs)
	}
}

func TestServiceShouldMapMetricKeysWhenAggregatingAllRegisteredTypes(t *testing.T) {
	querier := &recordingQuerier{result: 3, present: true}
	registry := testRegistry()
	service := aggregation.NewService(testRegistryProvider(registry), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to, "Europe/Moscow")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot.Timezone != "Europe/Moscow" {
		t.Fatalf("expected Europe/Moscow timezone, got %s", snapshot.Timezone)
	}
	if snapshot.Metrics["listingsPublished"].Value == nil || *snapshot.Metrics["listingsPublished"].Value != 3 {
		t.Fatalf("expected listingsPublished=3, got %v", snapshot.Metrics["listingsPublished"].Value)
	}
	if len(snapshot.Metrics) != len(registry.All()) {
		t.Fatalf("expected %d metrics, got %d", len(registry.All()), len(snapshot.Metrics))
	}
}

func TestServiceShouldReturnNullWhenSparseMetricIsAbsent(t *testing.T) {
	querier := &recordingQuerier{present: false}
	registry := testRegistry()
	service := aggregation.NewService(testRegistryProvider(registry), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to, "UTC")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot.Metrics["activeListings"].Value != nil {
		t.Fatalf("expected null activeListings value, got %v", snapshot.Metrics["activeListings"].Value)
	}
	if snapshot.Metrics["activeListings"].Percentile != nil {
		t.Fatalf("expected null activeListings percentile, got %v", snapshot.Metrics["activeListings"].Percentile)
	}
	if snapshot.Metrics["activeListings"].Share != nil {
		t.Fatalf("expected null activeListings share, got %v", snapshot.Metrics["activeListings"].Share)
	}
	if snapshot.Metrics["listingsPublished"].Value == nil || *snapshot.Metrics["listingsPublished"].Value != 0 {
		t.Fatalf("expected listingsPublished=0, got %v", snapshot.Metrics["listingsPublished"].Value)
	}
}

func TestServiceShouldFailWhenFromIsNotBeforeTo(t *testing.T) {
	querier := &recordingQuerier{}
	service := aggregation.NewService(testRegistryProvider(testRegistry()), querier)
	from := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := service.Metrics(context.Background(), 42, from, to, "UTC")

	if !apperr.IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestServiceShouldUseUTCWhenTimezoneIsEmpty(t *testing.T) {
	querier := &recordingQuerier{result: 1, present: true}
	service := aggregation.NewService(testRegistryProvider(testRegistry()), querier)
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	snapshot, err := service.Metrics(context.Background(), 42, from, to, "  ")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot.Timezone != "UTC" {
		t.Fatalf("expected UTC timezone, got %s", snapshot.Timezone)
	}
}
