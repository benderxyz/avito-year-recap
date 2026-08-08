package aggregation

import (
	"context"
	"fmt"
	"time"

	"analytics-service/internal/events"
)

type UniqueAggregator struct {
	db           FloatQuerier
	mode         events.UniqueMode
	payloadField string
}

func NewUniqueAggregator(db FloatQuerier, mode events.UniqueMode, payloadField string) *UniqueAggregator {
	if mode == "" {
		mode = events.UniqueModePayload
	}

	return &UniqueAggregator{
		db:           db,
		mode:         mode,
		payloadField: sanitizePayloadField(payloadField),
	}
}

func (a *UniqueAggregator) Aggregate(ctx context.Context, req Request) (Result, error) {
	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}
	if _, err := time.LoadLocation(timezone); err != nil {
		timezone = "UTC"
	}

	var (
		query string
		args  []any
	)

	switch a.mode {
	case events.UniqueModeDay:
		query = `
			SELECT toFloat64(uniqExact(toDate(toTimeZone(occurred_at, ?))))
			FROM events
			WHERE user_id = ?
			  AND event_type = ?
			  AND event_category = 'unique'
			  AND occurred_at >= ?
			  AND occurred_at < ?
		`
		args = []any{timezone, req.UserID, req.EventType, req.From, req.To}
	default:
		query = fmt.Sprintf(`
			SELECT toFloat64(uniqExact(JSONExtractString(payload, '%s')))
			FROM events
			WHERE user_id = ?
			  AND event_type = ?
			  AND event_category = 'unique'
			  AND JSONExtractString(payload, '%s') != ''
			  AND occurred_at >= ?
			  AND occurred_at < ?
		`, a.payloadField, a.payloadField)
		args = []any{req.UserID, req.EventType, req.From, req.To}
	}

	value, present, err := a.db.QueryFloat64(ctx, query, args...)
	if err != nil {
		return Result{}, fmt.Errorf("unique aggregate: %w", err)
	}
	if !present {
		return Result{Value: 0, Present: true}, nil
	}
	return Result{Value: value, Present: true}, nil
}

func (a *UniqueAggregator) PayloadField() string {
	return a.payloadField
}

func (a *UniqueAggregator) Mode() events.UniqueMode {
	return a.mode
}
