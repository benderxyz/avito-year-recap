package aggregation

import (
	"context"
	"time"
)

type FloatQuerier interface {
	QueryFloat64(ctx context.Context, query string, args ...any) (float64, bool, error)
}

type Request struct {
	UserID    uint64
	EventType string
	From      time.Time
	To        time.Time
	Timezone  string
}

type Aggregator interface {
	Aggregate(ctx context.Context, req Request) (Result, error)
}

type TimezoneResolver interface {
	Timezone(ctx context.Context, userID uint64) (string, error)
}
