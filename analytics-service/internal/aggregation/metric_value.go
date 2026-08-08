package aggregation

import "math"

type MetricValue struct {
	Value      *float64 `json:"value"`
	Percentile *float64 `json:"percentile"`
	Share      *float64 `json:"share"`
}

func nullMetricValue() MetricValue {
	return MetricValue{}
}

func metricValueWithValue(value float64) MetricValue {
	rounded := roundMetricFloat(value)
	return MetricValue{Value: &rounded}
}

func roundMetricFloatPtr(value float64) *float64 {
	rounded := roundMetricFloat(value)
	return &rounded
}

func roundMetricFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

func calculateShare(userValue, globalValue float64) *float64 {
	if globalValue <= 0 {
		return nil
	}
	share := roundMetricFloat(userValue / globalValue * 100)
	return &share
}
