package aggregation

import (
	"context"
	"fmt"

	"analytics-service/internal/events"
)

func enrichMetric(
	ctx context.Context,
	stats *globalStats,
	cfg events.CategoryConfig,
	eventType string,
	req Request,
	user Result,
) (MetricValue, error) {
	globalReq := GlobalRequest{
		EventType: eventType,
		From:      req.From,
		To:        req.To,
		Timezone:  req.Timezone,
	}

	switch cfg.Category {
	case events.CategoryCounter:
		return enrichCounterMetric(ctx, stats, globalReq, user)
	case events.CategoryUnique:
		return enrichUniqueMetric(ctx, stats, cfg, globalReq, user)
	case events.CategoryGauge:
		return enrichSparseLowerMetric(ctx, globalReq, user, stats.gaugeTotals)
	case events.CategoryMilestone:
		return enrichSparseHigherMetric(ctx, globalReq, user, stats.milestoneTotals)
	case events.CategoryInterval:
		return enrichSparseHigherMetric(ctx, globalReq, user, stats.intervalTotals)
	default:
		return MetricValue{}, fmt.Errorf("unsupported category for enrichment: %s", cfg.Category)
	}
}

func enrichCounterMetric(
	ctx context.Context,
	stats *globalStats,
	globalReq GlobalRequest,
	user Result,
) (MetricValue, error) {
	value := user.Value
	metric := metricValueWithValue(value)

	global, err := stats.counterTotal(ctx, globalReq)
	if err != nil {
		return MetricValue{}, err
	}
	metric.Share = calculateShare(value, global.Value)

	totals, err := stats.counterTotals(ctx, globalReq)
	if err != nil {
		return MetricValue{}, err
	}
	percentile := lowerPercentile(totals, value)
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}

func enrichUniqueMetric(
	ctx context.Context,
	stats *globalStats,
	cfg events.CategoryConfig,
	globalReq GlobalRequest,
	user Result,
) (MetricValue, error) {
	value := user.Value
	metric := metricValueWithValue(value)

	global, err := stats.uniqueTotal(ctx, globalReq, cfg)
	if err != nil {
		return MetricValue{}, err
	}
	metric.Share = calculateShare(value, global.Value)

	totals, err := stats.uniqueTotals(ctx, globalReq, cfg)
	if err != nil {
		return MetricValue{}, err
	}
	percentile := lowerPercentile(totals, value)
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}

func enrichSparseLowerMetric(
	ctx context.Context,
	globalReq GlobalRequest,
	user Result,
	totalsFn func(context.Context, GlobalRequest) ([]float64, error),
) (MetricValue, error) {
	if !user.Present {
		return nullMetricValue(), nil
	}

	value := user.Value
	metric := metricValueWithValue(value)

	totals, err := totalsFn(ctx, globalReq)
	if err != nil {
		return MetricValue{}, err
	}
	percentile := lowerPercentile(totals, value)
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}

func enrichSparseHigherMetric(
	ctx context.Context,
	globalReq GlobalRequest,
	user Result,
	totalsFn func(context.Context, GlobalRequest) ([]float64, error),
) (MetricValue, error) {
	if !user.Present {
		return nullMetricValue(), nil
	}

	value := user.Value
	metric := metricValueWithValue(value)

	totals, err := totalsFn(ctx, globalReq)
	if err != nil {
		return MetricValue{}, err
	}
	percentile := higherPercentile(totals, value)
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}
