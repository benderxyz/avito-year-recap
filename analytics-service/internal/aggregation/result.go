package aggregation

import "time"

type Result struct {
	Value   float64
	Present bool
}

type MetricsSnapshot struct {
	Metrics  map[string]MetricValue
	Timezone string
	From     time.Time
	To       time.Time
}
