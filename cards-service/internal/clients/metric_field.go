package clients

type MetricField struct {
	Value      *float64 `json:"value"`
	Percentile *float64 `json:"percentile"`
	Share      *float64 `json:"share"`
}
