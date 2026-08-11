package models

type MetricDefinition struct {
	Key           string
	ValueType     MetricType
	Currency      Currency
	IsPublic      bool
	PercentileKey string
	SourceKey     string
	SourceField   MetricSourceField
	IncludeInLLM  bool
}
