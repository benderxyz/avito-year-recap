package events_test

import "analytics-service/internal/events"

func testRegistry() *events.Registry {
	return events.NewRegistry(map[string]events.CategoryConfig{
		"item_published":       {Category: events.CategoryCounter, MetricKey: "listingsPublished"},
		"item_view_received":   {Category: events.CategoryCounter, MetricKey: "viewsTotal"},
		"first_item_published": {Category: events.CategoryMilestone, MetricKey: "firstListingAt"},
		"chat_reply":           {Category: events.CategoryInterval, MetricKey: "avgReplySeconds"},
		"category_opened": {
			Category:    events.CategoryUnique,
			MetricKey:   "categoriesTried",
			UniqueMode:  events.UniqueModePayload,
			UniqueField: "category",
		},
		"day_active": {
			Category:   events.CategoryUnique,
			MetricKey:  "daysActive",
			UniqueMode: events.UniqueModeDay,
		},
	})
}
