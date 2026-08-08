package models

type MetricKey string

const (
	MetricListingsPublished   MetricKey = "listingsPublished"
	MetricListingsPercentile  MetricKey = "listingsPercentile"
	MetricViewsTotal          MetricKey = "viewsTotal"
	MetricViewsPercentile     MetricKey = "viewsPercentile"
	MetricFavoritesReceived   MetricKey = "favoritesReceived"
	MetricFavoritesPercentile MetricKey = "favoritesPercentile"
	MetricMessagesSent        MetricKey = "messagesSent"
	MetricMessagesPercentile  MetricKey = "messagesPercentile"
	MetricDealsClosed         MetricKey = "dealsClosed"
	MetricDealsPercentile     MetricKey = "dealsPercentile"
	MetricMoneyEarned         MetricKey = "moneyEarned"
	MetricMoneySaved          MetricKey = "moneySaved"
	MetricDaysActive          MetricKey = "daysActive"
	MetricPeakDayViews        MetricKey = "peakDayViews"
	MetricSearchQueries       MetricKey = "searchQueries"
	MetricCategoriesTried     MetricKey = "categoriesTried"
	MetricDeliveryOrders      MetricKey = "deliveryOrders"
	MetricActiveListings      MetricKey = "activeListings"
	MetricSellerRating        MetricKey = "sellerRating"
	MetricAvgReplySeconds     MetricKey = "avgReplySeconds"
	MetricFirstListingAt      MetricKey = "firstListingAt"
	MetricFirstDealAt         MetricKey = "firstDealAt"
)

type MetricType string

const (
	MetricTypeNumber     MetricType = "number"
	MetricTypeMoney      MetricType = "money"
	MetricTypePercentile MetricType = "percentile"
)

type Currency string

const (
	CurrencyRUB Currency = "RUB"
)
