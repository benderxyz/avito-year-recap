package aggregation

import (
	"context"
	"fmt"
)

type GaugeAggregator struct {
	db FloatQuerier
}

func NewGaugeAggregator(db FloatQuerier) *GaugeAggregator {
	return &GaugeAggregator{db: db}
}

func (a *GaugeAggregator) Aggregate(ctx context.Context, req Request) (Result, error) {
	const query = `
		SELECT if(count() = 0, CAST(NULL AS Nullable(Float64)), argMax(value, occurred_at))
		FROM events
		WHERE user_id = ?
		  AND event_type = ?
		  AND event_category = 'gauge'
		  AND occurred_at >= ?
		  AND occurred_at < ?
	`

	value, present, err := a.db.QueryFloat64(ctx, query, req.UserID, req.EventType, req.From, req.To)
	if err != nil {
		return Result{}, fmt.Errorf("gauge aggregate: %w", err)
	}
	return Result{Value: value, Present: present}, nil
}
