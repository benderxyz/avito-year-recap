package clients

import "math"

func ParseMetrics(fields map[string]MetricField) Metrics {
	return Metrics{
		ListingsPublished:   intValue(fields, "listingsPublished"),
		ListingsPercentile:  percentileValue(fields, "listingsPublished"),
		ViewsTotal:          intValue(fields, "viewsTotal"),
		ViewsPercentile:     percentileValue(fields, "viewsTotal"),
		FavoritesReceived:   intValue(fields, "favoritesReceived"),
		FavoritesPercentile: percentileValue(fields, "favoritesReceived"),
		MessagesSent:        intValue(fields, "messagesSent"),
		MessagesPercentile:  percentileValue(fields, "messagesSent"),
		DealsClosed:         intValue(fields, "dealsClosed"),
		DealsPercentile:     percentileValue(fields, "dealsClosed"),
		MoneyEarned:         int64Value(fields, "moneyEarned"),
		MoneySaved:          int64Value(fields, "moneySaved"),
		DaysActive:          intValue(fields, "daysActive"),
		PeakDayViews:        intValue(fields, "peakDayViews"),
		SearchQueries:       intValue(fields, "searchQueries"),
		CategoriesTried:     intValue(fields, "categoriesTried"),
		DeliveryOrders:      intValue(fields, "deliveryOrders"),
		ActiveListings:      intValue(fields, "activeListings"),
		AvgReplySeconds:     floatValue(fields, "avgReplySeconds"),
		SellerRating:        floatValue(fields, "sellerRating"),
		FirstDealAt:         int64Value(fields, "firstDealAt"),
		FirstListingAt:      int64Value(fields, "firstListingAt"),
	}
}

func intValue(fields map[string]MetricField, key string) int {
	value, ok := scalarValue(fields, key)
	if !ok {
		return 0
	}
	return int(value)
}

func int64Value(fields map[string]MetricField, key string) int64 {
	value, ok := scalarValue(fields, key)
	if !ok {
		return 0
	}
	return int64(value)
}

func floatValue(fields map[string]MetricField, key string) float64 {
	value, ok := scalarValue(fields, key)
	if !ok {
		return 0
	}
	return value
}

func percentileValue(fields map[string]MetricField, key string) *float64 {
	field, ok := fields[key]
	if !ok || field.Percentile == nil {
		return nil
	}
	rounded := roundToTwoDecimals(*field.Percentile)
	return &rounded
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

func scalarValue(fields map[string]MetricField, key string) (float64, bool) {
	field, ok := fields[key]
	if !ok || field.Value == nil {
		return 0, false
	}
	return *field.Value, true
}
