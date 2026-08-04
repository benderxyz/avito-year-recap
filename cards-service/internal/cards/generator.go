package cards

import (
	"fmt"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

func Generate(metrics clients.Metrics) []models.Card {
	var cards []models.Card

	if metrics.ViewsRealty > 50 {
		cards = append(cards, models.Card{
			Type:  "achievement",
			Title: "Искатель жилья",
			Text:  fmt.Sprintf("Ты посмотрела %d объявлений о недвижимости", metrics.ViewsRealty),
		})
	}

	if metrics.ViewsAuto > 30 {
		cards = append(cards, models.Card{
			Type:  "achievement",
			Title: "Автоман",
			Text:  fmt.Sprintf("Ты посмотрела %d автомобилей", metrics.ViewsAuto),
		})
	}

	if metrics.Purchases >= 3 {
		cards = append(cards, models.Card{
			Type:  "achievement",
			Title: "Уже с опытом покупок",
			Text:  fmt.Sprintf("За год ты совершила %d покупок", metrics.Purchases),
		})
	}

	if metrics.ViewsRealty < 10 && metrics.ViewsAuto < 10 && metrics.Purchases == 0 {
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
