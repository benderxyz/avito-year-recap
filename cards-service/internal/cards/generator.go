package cards

import (
	"fmt"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

func Generate(metrics clients.Metrics) []models.Card {
	var cards []models.Card

	if metrics.ViewsTotal > 500 {
		cards = append(cards, models.Card{
			Type:  "achievement",
			Title: "Много внимания",
			Text:  fmt.Sprintf("Твои объявления набрали %d просмотров", metrics.ViewsTotal),
		})
	}

	if metrics.FavoritesReceived > 20 {
		cards = append(cards, models.Card{
			Type:  "achievement",
			Title: "В избранном",
			Text:  fmt.Sprintf("Твои объявления добавили в избранное %d раз", metrics.FavoritesReceived),
		})
	}

	if metrics.DealsClosed >= 3 {
		cards = append(cards, models.Card{
			Type:  "achievement",
			Title: "Сделки пошли",
			Text:  fmt.Sprintf("За год у тебя было %d успешных сделок", metrics.DealsClosed),
		})
	}

	if metrics.DaysActive < 10 && metrics.ViewsTotal == 0 && metrics.DealsClosed == 0 {
		cards = append(cards, models.Card{
			Type:  "insight",
			Title: "Тихий наблюдатель",
			Text:  "В этом году ты почти не пользовалась Авито. Может, самое время начать?",
		})
	}

	cards = append(cards, models.Card{
		Type:        "recommendation",
		Title:       "Что дальше?",
		Text:        "Посмотри, что нового появилось в интересных тебе категориях",
		Action:      "open_home",
		ActionValue: "",
	})

	return cards
}
