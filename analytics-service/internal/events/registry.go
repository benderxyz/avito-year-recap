package events

import "fmt"

var defaultRegistry = map[string]CategoryConfig{
	"item_published": {
		Category:  CategoryCounter,
		MetricKey: "listingsPublished",
	},
	"item_view_received": {
		Category:  CategoryCounter,
		MetricKey: "viewsTotal",
	},
	"item_favorited": {
		Category:  CategoryCounter,
		MetricKey: "favoritesReceived",
	},
	"chat_message_sent": {
		Category:  CategoryCounter,
		MetricKey: "messagesSent",
	},
	"deal_completed": {
		Category:  CategoryCounter,
		MetricKey: "dealsClosed",
	},
	"sale_proceeds": {
		Category:  CategoryCounter,
		MetricKey: "moneyEarned",
	},
	"buyer_saving": {
		Category:  CategoryCounter,
		MetricKey: "moneySaved",
	},
	"delivery_order_created": {
		Category:  CategoryCounter,
		MetricKey: "deliveryOrders",
	},
	"search_performed": {
		Category:  CategoryCounter,
		MetricKey: "searchQueries",
	},
	"active_items_count": {
		Category:  CategoryGauge,
		MetricKey: "activeListings",
	},
	"seller_rating": {
		Category:  CategoryGauge,
		MetricKey: "sellerRating",
	},
	"peak_day_views": {
		Category:  CategoryGauge,
		MetricKey: "peakDayViews",
	},
	"first_item_published": {
		Category:  CategoryMilestone,
		MetricKey: "firstListingAt",
	},
	"first_deal_completed": {
		Category:  CategoryMilestone,
		MetricKey: "firstDealAt",
	},
	"category_opened": {
		Category:    CategoryUnique,
		UniqueMode:  UniqueModePayload,
		UniqueField: "category",
		MetricKey:   "categoriesTried",
	},
	"day_active": {
		Category:   CategoryUnique,
		UniqueMode: UniqueModeDay,
		MetricKey:  "daysActive",
	},
	"chat_reply": {
		Category:  CategoryInterval,
		MetricKey: "avgReplySeconds",
	},
}

type Registry struct {
	byType map[string]CategoryConfig
}

func NewRegistry() *Registry {
	copied := make(map[string]CategoryConfig, len(defaultRegistry))
	for key, value := range defaultRegistry {
		copied[key] = value
	}
	return &Registry{byType: copied}
}

func (r *Registry) Get(eventType string) (CategoryConfig, error) {
	cfg, ok := r.byType[eventType]
	if !ok {
		return CategoryConfig{}, fmt.Errorf("unknown event type: %s", eventType)
	}
	return cfg, nil
}

func (r *Registry) All() map[string]CategoryConfig {
	copied := make(map[string]CategoryConfig, len(r.byType))
	for key, value := range r.byType {
		copied[key] = value
	}
	return copied
}

func (r *Registry) Has(eventType string) bool {
	_, ok := r.byType[eventType]
	return ok
}
