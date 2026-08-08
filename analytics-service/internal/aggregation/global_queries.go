package aggregation

import (
	"context"
	"fmt"
	"time"

	"analytics-service/internal/events"
)

type GlobalRequest struct {
	EventType string
	From      time.Time
	To        time.Time
	Timezone  string
}

func aggregateGlobalCounter(ctx context.Context, db FloatQuerier, req GlobalRequest) (Result, error) {
	const query = `
		SELECT sum(value)
		FROM events
		WHERE event_type = ?
		  AND event_category = 'counter'
		  AND occurred_at >= ?
		  AND occurred_at < ?
	`

	value, present, err := db.QueryFloat64(ctx, query, req.EventType, req.From, req.To)
	if err != nil {
		return Result{}, fmt.Errorf("global counter aggregate: %w", err)
	}
	if !present {
		return Result{Value: 0, Present: true}, nil
	}
	return Result{Value: value, Present: true}, nil
}

func aggregateGlobalUnique(
	ctx context.Context,
	db FloatQuerier,
	req GlobalRequest,
	cfg events.CategoryConfig,
) (Result, error) {
	subquery, args, err := uniquePerUserSubquery(cfg, req)
	if err != nil {
		return Result{}, err
	}

	query := fmt.Sprintf(`
		SELECT sum(user_total)
		FROM (%s)
	`, subquery)

	value, present, err := db.QueryFloat64(ctx, query, args...)
	if err != nil {
		return Result{}, fmt.Errorf("global unique aggregate: %w", err)
	}
	if !present {
		return Result{Value: 0, Present: true}, nil
	}
	return Result{Value: value, Present: true}, nil
}

func counterPerUserSubquery(req GlobalRequest) (string, []any) {
	const query = `
		SELECT sum(value) AS total
		FROM events
		WHERE event_type = ?
		  AND event_category = 'counter'
		  AND occurred_at >= ?
		  AND occurred_at < ?
		GROUP BY user_id
	`
	return query, []any{req.EventType, req.From, req.To}
}

func gaugePerUserSubquery(req GlobalRequest) (string, []any) {
	const query = `
		SELECT argMax(value, occurred_at) AS total
		FROM events
		WHERE event_type = ?
		  AND event_category = 'gauge'
		  AND occurred_at >= ?
		  AND occurred_at < ?
		GROUP BY user_id
	`
	return query, []any{req.EventType, req.From, req.To}
}

func milestonePerUserSubquery(req GlobalRequest) (string, []any) {
	const query = `
		SELECT toFloat64(toUnixTimestamp(min(occurred_at))) AS total
		FROM events
		WHERE event_type = ?
		  AND event_category = 'milestone'
		  AND occurred_at >= ?
		  AND occurred_at < ?
		GROUP BY user_id
	`
	return query, []any{req.EventType, req.From, req.To}
}

func intervalPerUserSubquery(req GlobalRequest) (string, []any) {
	const query = `
		SELECT avg(toUnixTimestamp(occurred_at) - value) AS total
		FROM events
		WHERE event_type = ?
		  AND event_category = 'interval'
		  AND occurred_at >= ?
		  AND occurred_at < ?
		GROUP BY user_id
	`
	return query, []any{req.EventType, req.From, req.To}
}

func uniquePerUserSubquery(cfg events.CategoryConfig, req GlobalRequest) (string, []any, error) {
	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	switch cfg.UniqueMode {
	case events.UniqueModeDay:
		const query = `
			SELECT toFloat64(uniqExact(toDate(toTimeZone(occurred_at, ?)))) AS user_total
			FROM events
			WHERE event_type = ?
			  AND event_category = 'unique'
			  AND occurred_at >= ?
			  AND occurred_at < ?
			GROUP BY user_id
		`
		return query, []any{timezone, req.EventType, req.From, req.To}, nil
	default:
		field := sanitizePayloadField(cfg.UniqueField)
		query := fmt.Sprintf(`
			SELECT toFloat64(uniqExact(JSONExtractString(payload, '%s'))) AS user_total
			FROM events
			WHERE event_type = ?
			  AND event_category = 'unique'
			  AND JSONExtractString(payload, '%s') != ''
			  AND occurred_at >= ?
			  AND occurred_at < ?
			GROUP BY user_id
		`, field, field)
		return query, []any{req.EventType, req.From, req.To}, nil
	}
}

func uniquePerUserTotalsSubquery(cfg events.CategoryConfig, req GlobalRequest) (string, []any, error) {
	subquery, args, err := uniquePerUserSubquery(cfg, req)
	if err != nil {
		return "", nil, err
	}
	query := fmt.Sprintf(`
		SELECT user_total AS total
		FROM (%s)
	`, subquery)
	return query, args, nil
}
