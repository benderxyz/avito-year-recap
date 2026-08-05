package aggregation

import (
	"context"
	"fmt"
)

type IntervalAggregator struct {
	db FloatQuerier
}

func NewIntervalAggregator(db FloatQuerier) *IntervalAggregator {
	return &IntervalAggregator{db: db}
}

func (a *IntervalAggregator) Aggregate(ctx context.Context, req Request) (Result, error) {
	const query = `
		SELECT if(count() = 0, CAST(NULL AS Nullable(Float64)), avg(toUnixTimestamp(occurred_at) - value))
		FROM events
		WHERE user_id = ?
		  AND event_type = ?
		  AND event_category = 'interval'
		  AND occurred_at >= ?
		  AND occurred_at < ?
	`

	value, present, err := a.db.QueryFloat64(ctx, query, req.UserID, req.EventType, req.From, req.To)
	if err != nil {
		return Result{}, fmt.Errorf("interval aggregate: %w", err)
	}
	return Result{Value: value, Present: present}, nil
}
