package aggregation_test

import (
	"context"

	"analytics-service/internal/events"
)

type staticRegistryProvider struct {
	registry *events.Registry
}

func (p staticRegistryProvider) Get(context.Context) (*events.Registry, error) {
	return p.registry, nil
}

func testRegistry() *events.Registry {
	return events.NewRegistry(map[string]events.CategoryConfig{
		"item_published":       {Category: events.CategoryCounter, MetricKey: "listingsPublished"},
		"active_items_count":   {Category: events.CategoryGauge, MetricKey: "activeListings"},
		"first_item_published": {Category: events.CategoryMilestone, MetricKey: "firstListingAt"},
		"chat_reply":           {Category: events.CategoryInterval, MetricKey: "avgReplySeconds"},
		"day_active": {
			Category:   events.CategoryUnique,
			MetricKey:  "daysActive",
			UniqueMode: events.UniqueModeDay,
		},
	})
}

func testRegistryProvider(registry *events.Registry) staticRegistryProvider {
	return staticRegistryProvider{registry: registry}
}
