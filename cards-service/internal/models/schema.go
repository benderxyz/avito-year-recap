package models

type MetricKey string

type MetricType string

const (
	MetricTypeNumber     MetricType = "number"
	MetricTypeMoney      MetricType = "money"
	MetricTypePercentile MetricType = "percentile"
	MetricTypeRatio      MetricType = "ratio"
	MetricTypeString     MetricType = "string"
	MetricTypeDate       MetricType = "date"
)

type MetricSourceField string

const (
	MetricSourceValue      MetricSourceField = "value"
	MetricSourcePercentile MetricSourceField = "percentile"
	MetricSourceShare      MetricSourceField = "share"
)

type Currency string

const (
	CurrencyRUB Currency = "RUB"
)
