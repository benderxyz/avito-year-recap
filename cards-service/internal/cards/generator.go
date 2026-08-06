package cards

import (
	"cards-service/internal/clients"
	"cards-service/internal/models"
	"fmt"
)

/*
	type Metrics struct {
		ActiveListings    int     `json:"activeListings"`
		AvgReplySeconds   float64 `json:"avgReplySeconds"`
		CategoriesTried   int     `json:"categoriesTried"`
		DaysActive        int     `json:"daysActive"`
		DealsClosed       int     `json:"dealsClosed"`
		DeliveryOrders    int     `json:"deliveryOrders"`
		FavoritesReceived int     `json:"favoritesReceived"`
		FirstDealAt       int64   `json:"firstDealAt"`
		FirstListingAt    int64   `json:"firstListingAt"`
		ListingsPublished int     `json:"listingsPublished"`
		MessagesSent      int     `json:"messagesSent"`
		MoneyEarned       int64   `json:"moneyEarned"`
		MoneySaved        int64   `json:"moneySaved"`
		PeakDayViews      int     `json:"peakDayViews"`
		SearchQueries     int     `json:"searchQueries"`
		SellerRating      float64 `json:"sellerRating"`
		ViewsTotal        int     `json:"viewsTotal"`
	}
*/
func Generate(metrics clients.Metrics) []models.Badge {
	var badges []models.Badge

	if messagesSent := metrics.MessagesSent; messagesSent > 1000 {
		badges = append(badges, models.Badge{
			ID:          "messages_sent_1000",
			Title:       "Разговорчивый",
			Description: fmt.Sprintf("Вы отправили %d сообщений!", messagesSent),
		})
	}

	if metrics != (clients.Metrics{}) {
		badges = append(badges, models.Badge{
			ID:          "active_user",
			Title:       "Активный пользователь",
			Description: "Вы заходили на нашу площадку в этом году!",
		})
	}

	fmt.Println(badges)

	return badges
}
