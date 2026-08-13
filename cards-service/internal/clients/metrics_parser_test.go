package clients

import "testing"

func TestParseMetricsShouldKeepAnyKeyWhenAnalyticsReturnsMetricObjects(t *testing.T) {
	metrics := ParseMetrics(map[string]MetricField{
		"listingsPublished": {
			Value:      floatPtr(47),
			Percentile: floatPtr(88.4),
			Share:      floatPtr(4.7),
		},
		"brandNewMetric": {
			Value: floatPtr(3),
		},
		"activeListings": {
			Value:      nil,
			Percentile: nil,
			Share:      nil,
		},
	})

	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}
	if metrics["listingsPublished"].Value == nil || *metrics["listingsPublished"].Value != 47 {
		t.Fatalf("expected listingsPublished 47, got %v", metrics["listingsPublished"].Value)
	}
	if metrics["brandNewMetric"].Value == nil || *metrics["brandNewMetric"].Value != 3 {
		t.Fatalf("expected unknown key to pass through, got %v", metrics["brandNewMetric"].Value)
	}
	if metrics["activeListings"].Value != nil {
		t.Fatalf("expected nil value to stay nil, got %v", metrics["activeListings"].Value)
	}
}

func TestParseMetricsShouldRoundPercentileAndShareWhenValuesHaveExtraPrecision(t *testing.T) {
	metrics := ParseMetrics(map[string]MetricField{
		"viewsTotal": {
			Value:      floatPtr(12840),
			Percentile: floatPtr(92.1449),
			Share:      floatPtr(4.7891),
		},
	})

	sample := metrics["viewsTotal"]
	if sample.Percentile == nil || *sample.Percentile != 92.14 {
		t.Fatalf("expected percentile 92.14, got %v", sample.Percentile)
	}
	if sample.Share == nil || *sample.Share != 4.79 {
		t.Fatalf("expected share 4.79, got %v", sample.Share)
	}
}

func floatPtr(value float64) *float64 {
	return &value
}
