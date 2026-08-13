package clients

import "math"

func ParseMetrics(fields map[string]MetricField) Metrics {
	metrics := make(Metrics, len(fields))
	for key, field := range fields {
		metrics[key] = MetricSample{
			Value:      field.Value,
			Percentile: roundPtr(field.Percentile),
			Share:      roundPtr(field.Share),
		}
	}
	return metrics
}

func roundPtr(value *float64) *float64 {
	if value == nil {
		return nil
	}
	rounded := roundToTwoDecimals(*value)
	return &rounded
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
