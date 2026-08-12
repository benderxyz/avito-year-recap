package models

type MetricDefinition struct {
	Key                     string
	ValueType               MetricType
	Currency                Currency
	IsPublic                bool
	PercentileKey           string
	ComparisonMinPercentile float64
	SourceKey               string
	SourceField             MetricSourceField
	IncludeInLLM            bool
}
