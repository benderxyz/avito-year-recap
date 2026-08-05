package aggregation

import (
	"fmt"

	"analytics-service/internal/events"
)

func NewAggregator(db FloatQuerier, cfg events.CategoryConfig) (Aggregator, error) {
	switch cfg.Category {
	case events.CategoryCounter:
		return NewCounterAggregator(db), nil
	case events.CategoryInterval:
		return NewIntervalAggregator(db), nil
	case events.CategoryGauge:
		return NewGaugeAggregator(db), nil
	case events.CategoryMilestone:
		return NewMilestoneAggregator(db), nil
	case events.CategoryUnique:
		return NewUniqueAggregator(db, cfg.UniqueMode, cfg.UniqueField), nil
	default:
		return nil, fmt.Errorf("unsupported event category: %s", cfg.Category)
	}
}
