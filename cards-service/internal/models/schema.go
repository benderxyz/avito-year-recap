package models

type MetricKey string

const (
	MetricListingsPublished MetricKey = "listingsPublished"
	MetricViewsTotal        MetricKey = "viewsTotal"
	MetricFavoritesReceived MetricKey = "favoritesReceived"
	MetricMessagesSent      MetricKey = "messagesSent"
	MetricDealsClosed       MetricKey = "dealsClosed"
	MetricMoneyEarned       MetricKey = "moneyEarned"
	MetricMoneySaved        MetricKey = "moneySaved"
	MetricDaysActive        MetricKey = "daysActive"
	MetricPeakDayViews      MetricKey = "peakDayViews"
	MetricSearchQueries     MetricKey = "searchQueries"
	MetricCategoriesTried   MetricKey = "categoriesTried"
	MetricDeliveryOrders    MetricKey = "deliveryOrders"
)

type MetricType string

const (
	MetricTypeNumber MetricType = "number"
	MetricTypeMoney  MetricType = "money"
)

type Currency string

const (
	CurrencyRUB Currency = "RUB"
)
