package aggregation

import (
	"context"
	"fmt"
)

type MilestoneAggregator struct {
	db FloatQuerier
}

func NewMilestoneAggregator(db FloatQuerier) *MilestoneAggregator {
	return &MilestoneAggregator{db: db}
}

func (a *MilestoneAggregator) Aggregate(ctx context.Context, req Request) (Result, error) {
	const query = `
		SELECT if(
			count() = 0,
			CAST(NULL AS Nullable(Float64)),
			toFloat64(toUnixTimestamp(min(occurred_at)))
		)
		FROM events
		WHERE user_id = ?
		  AND event_type = ?
		  AND event_category = 'milestone'
		  AND occurred_at >= ?
		  AND occurred_at < ?
	`

	value, present, err := a.db.QueryFloat64(ctx, query, req.UserID, req.EventType, req.From, req.To)
	if err != nil {
		return Result{}, fmt.Errorf("milestone aggregate: %w", err)
	}
	return Result{Value: value, Present: present}, nil
}
