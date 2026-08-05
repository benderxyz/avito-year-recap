package aggregation

import (
	"context"
	"fmt"
	"time"

	"analytics-service/internal/apperr"
	"analytics-service/internal/events"
)

type Service struct {
	registry  *events.Registry
	db        FloatQuerier
	timezones TimezoneResolver
}

func NewService(registry *events.Registry, db FloatQuerier, timezones TimezoneResolver) *Service {
	return &Service{registry: registry, db: db, timezones: timezones}
}

func (s *Service) Metrics(
	ctx context.Context,
	userID uint64,
	from, to time.Time,
) (MetricsSnapshot, error) {
	if !from.Before(to) {
		return MetricsSnapshot{}, apperr.Validation("from must be before to")
	}

	timezone := "UTC"
	if s.timezones != nil {
		resolved, err := s.timezones.Timezone(ctx, userID)
		if err != nil {
			return MetricsSnapshot{}, fmt.Errorf("timezone: %w", err)
		}
		timezone = resolved
	}

	result := make(map[string]*float64)
	for eventType, cfg := range s.registry.All() {
		aggregator, err := NewAggregator(s.db, cfg)
		if err != nil {
			return MetricsSnapshot{}, fmt.Errorf("factory for %s: %w", eventType, err)
		}

		aggregate, err := aggregator.Aggregate(ctx, Request{
			UserID:    userID,
			EventType: eventType,
			From:      from,
			To:        to,
			Timezone:  timezone,
		})
		if err != nil {
			return MetricsSnapshot{}, fmt.Errorf("aggregate %s: %w", eventType, err)
		}

		key := cfg.MetricKey
		if key == "" {
			key = eventType
		}

		if aggregate.Present {
			value := aggregate.Value
			result[key] = &value
		} else {
			result[key] = nil
		}
	}

	return MetricsSnapshot{
		Metrics:  result,
		Timezone: timezone,
		From:     from,
		To:       to,
	}, nil
}
