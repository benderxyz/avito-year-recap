package aggregation

import (
	"context"
	"fmt"

	"analytics-service/internal/events"
)

func enrichMetric(
	ctx context.Context,
	db FloatQuerier,
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
		return enrichCounterMetric(ctx, db, globalReq, user)
	case events.CategoryUnique:
		return enrichUniqueMetric(ctx, db, cfg, globalReq, user)
	case events.CategoryGauge:
		return enrichSparseLowerMetric(ctx, db, globalReq, user, gaugePerUserSubquery)
	case events.CategoryMilestone:
		return enrichSparseHigherMetric(ctx, db, globalReq, user, milestonePerUserSubquery)
	case events.CategoryInterval:
		return enrichSparseHigherMetric(ctx, db, globalReq, user, intervalPerUserSubquery)
	default:
		return MetricValue{}, fmt.Errorf("unsupported category for enrichment: %s", cfg.Category)
	}
}

func enrichCounterMetric(
	ctx context.Context,
	db FloatQuerier,
	globalReq GlobalRequest,
	user Result,
) (MetricValue, error) {
	value := user.Value
	metric := metricValueWithValue(value)

	global, err := aggregateGlobalCounter(ctx, db, globalReq)
	if err != nil {
		return MetricValue{}, err
	}
	metric.Share = calculateShare(value, global.Value)

	subquery, args := counterPerUserSubquery(globalReq)
	percentile, err := aggregateLowerPercentile(ctx, db, subquery, args, value)
	if err != nil {
		return MetricValue{}, err
	}
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}

func enrichUniqueMetric(
	ctx context.Context,
	db FloatQuerier,
	cfg events.CategoryConfig,
	globalReq GlobalRequest,
	user Result,
) (MetricValue, error) {
	value := user.Value
	metric := metricValueWithValue(value)

	global, err := aggregateGlobalUnique(ctx, db, globalReq, cfg)
	if err != nil {
		return MetricValue{}, err
	}
	metric.Share = calculateShare(value, global.Value)

	subquery, args, err := uniquePerUserTotalsSubquery(cfg, globalReq)
	if err != nil {
		return MetricValue{}, err
	}
	percentile, err := aggregateLowerPercentile(ctx, db, subquery, args, value)
	if err != nil {
		return MetricValue{}, err
	}
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}

func enrichSparseLowerMetric(
	ctx context.Context,
	db FloatQuerier,
	globalReq GlobalRequest,
	user Result,
	subqueryFn func(GlobalRequest) (string, []any),
) (MetricValue, error) {
	if !user.Present {
		return nullMetricValue(), nil
	}

	value := user.Value
	metric := metricValueWithValue(value)

	subquery, args := subqueryFn(globalReq)
	percentile, err := aggregateLowerPercentile(ctx, db, subquery, args, value)
	if err != nil {
		return MetricValue{}, err
	}
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}

func enrichSparseHigherMetric(
	ctx context.Context,
	db FloatQuerier,
	globalReq GlobalRequest,
	user Result,
	subqueryFn func(GlobalRequest) (string, []any),
) (MetricValue, error) {
	if !user.Present {
		return nullMetricValue(), nil
	}

	value := user.Value
	metric := metricValueWithValue(value)

	subquery, args := subqueryFn(globalReq)
	percentile, err := aggregateHigherPercentile(ctx, db, subquery, args, value)
	if err != nil {
		return MetricValue{}, err
	}
	if percentile.Present {
		metric.Percentile = roundMetricFloatPtr(percentile.Value)
	}

	return metric, nil
}
