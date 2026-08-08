package clients

import "testing"

func TestParseMetricsShouldMapNestedFieldsWhenAnalyticsReturnsMetricObjects(t *testing.T) {
	metrics := ParseMetrics(map[string]MetricField{
		"listingsPublished": {
			Value:      floatPtr(47),
			Percentile: floatPtr(88.4),
			Share:      floatPtr(4.7),
		},
		"viewsTotal": {
			Value:      floatPtr(12840),
			Percentile: floatPtr(92.1),
		},
		"activeListings": {
			Value:      nil,
			Percentile: nil,
			Share:      nil,
		},
	})

	if metrics.ListingsPublished != 47 {
		t.Fatalf("expected listingsPublished 47, got %d", metrics.ListingsPublished)
	}
	if metrics.ListingsPercentile == nil || *metrics.ListingsPercentile != 88.4 {
		t.Fatalf("expected listings percentile 88.4, got %v", metrics.ListingsPercentile)
	}
	if metrics.ViewsTotal != 12840 {
		t.Fatalf("expected viewsTotal 12840, got %d", metrics.ViewsTotal)
	}
	if metrics.ActiveListings != 0 {
		t.Fatalf("expected activeListings 0 when value is null, got %d", metrics.ActiveListings)
	}
}

func floatPtr(value float64) *float64 {
	return &value
}
