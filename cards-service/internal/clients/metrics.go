package clients

import "cards-service/internal/models"

type MetricField struct {
	Value      *float64 `json:"value"`
	Percentile *float64 `json:"percentile"`
	Share      *float64 `json:"share"`
}

type MetricsResponse struct {
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	UserID   uint64                 `json:"user_id"`
	Timezone string                 `json:"timezone"`
	Metrics  map[string]MetricField `json:"metrics"`
}

type MetricSample struct {
	Value      *float64
	Percentile *float64
	Share      *float64
}

type Metrics map[string]MetricSample

func (m Metrics) Value(def models.MetricDefinition) (float64, bool) {
	sourceKey := def.SourceKey
	if sourceKey == "" {
		sourceKey = def.Key
	}

	sample, ok := m[sourceKey]
	if !ok {
		return 0, false
	}

	var value *float64
	switch def.SourceField {
	case models.MetricSourcePercentile:
		value = sample.Percentile
	case models.MetricSourceShare:
		value = sample.Share
	default:
		value = sample.Value
	}

	if value == nil {
		return 0, false
	}
	return *value, true
}
