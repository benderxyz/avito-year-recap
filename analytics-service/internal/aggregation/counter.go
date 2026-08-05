package aggregation

import (
	"context"
	"fmt"
)

type CounterAggregator struct {
	db FloatQuerier
}

func NewCounterAggregator(db FloatQuerier) *CounterAggregator {
	return &CounterAggregator{db: db}
}

func (a *CounterAggregator) Aggregate(ctx context.Context, req Request) (Result, error) {
	const query = `
		SELECT sum(value)
		FROM events
		WHERE user_id = ?
		  AND event_type = ?
		  AND event_category = 'counter'
		  AND occurred_at >= ?
		  AND occurred_at < ?
	`

	value, present, err := a.db.QueryFloat64(ctx, query, req.UserID, req.EventType, req.From, req.To)
	if err != nil {
		return Result{}, fmt.Errorf("counter aggregate: %w", err)
	}
	if !present {
		return Result{Value: 0, Present: true}, nil
	}
	return Result{Value: value, Present: true}, nil
}
