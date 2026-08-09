package aggregation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"analytics-service/internal/apperr"
	"analytics-service/internal/events"
)

type Service struct {
	registry *events.Registry
	db       FloatQuerier
	globals  *globalStats
}

func NewService(registry *events.Registry, db FloatQuerier) *Service {
	return &Service{
		registry: registry,
		db:       db,
		globals:  newGlobalStats(db, defaultGlobalCacheTTL),
	}
}

func (s *Service) Metrics(
	ctx context.Context,
	userID uint64,
	from, to time.Time,
	timezone string,
) (MetricsSnapshot, error) {
	if !from.Before(to) {
		return MetricsSnapshot{}, apperr.Validation("from must be before to")
	}

	timezone = strings.TrimSpace(timezone)
	if timezone == "" {
		timezone = "UTC"
	}

	result := make(map[string]MetricValue)
	for eventType, cfg := range s.registry.All() {
		aggregator, err := NewAggregator(s.db, cfg)
		if err != nil {
			return MetricsSnapshot{}, fmt.Errorf("factory for %s: %w", eventType, err)
		}

		req := Request{
			UserID:    userID,
			EventType: eventType,
			From:      from,
			To:        to,
			Timezone:  timezone,
		}

		aggregate, err := aggregator.Aggregate(ctx, req)
		if err != nil {
			return MetricsSnapshot{}, fmt.Errorf("aggregate %s: %w", eventType, err)
		}

		metric, err := enrichMetric(ctx, s.globals, cfg, eventType, req, aggregate)
		if err != nil {
			return MetricsSnapshot{}, fmt.Errorf("enrich %s: %w", eventType, err)
		}

		key := cfg.MetricKey
		if key == "" {
			key = eventType
		}
		result[key] = metric
	}

	return MetricsSnapshot{
		Metrics:  result,
		Timezone: timezone,
		From:     from,
		To:       to,
	}, nil
}
