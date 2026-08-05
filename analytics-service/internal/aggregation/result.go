package aggregation

import "time"

type Result struct {
	Value   float64
	Present bool
}

type MetricsSnapshot struct {
	Metrics  map[string]*float64
	Timezone string
	From     time.Time
	To       time.Time
}
