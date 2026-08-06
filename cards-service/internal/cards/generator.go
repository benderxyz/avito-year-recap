package cards

import (
	"fmt"
	"time"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

func BuildRecap(profile clients.Profile, year int, metrics clients.Metrics) models.RecapPayload {
	return models.RecapPayload{
		SchemaVersion: 1,
		Meta: models.Meta{
			Vertical: "marketplace",
			Year:     year,
			Locale:   "ru-RU",
			User: models.User{
				ID:          profile.ExternalID,
				DisplayName: profile.Username,
			},
			GeneratedAt: time.Now().UTC().Format(time.RFC3339), // в идеале здесь нужно передавать время генерации из сервиса аналитики
		},
		Metrics: buildMetrics(metrics),
		Badges:  buildBadges(metrics),
		Story:   buildStory(profile, year, metrics),
	}
}

func buildStory(profile clients.Profile, year int, m clients.Metrics) []map[string]any {
	name := profile.Username
	if name == "" {
		name = "вы"
	}

	story := make([]map[string]any, 0)

	story = append(story, map[string]any{
		"id":       "intro",
		"type":     "intro",
		"title":    fmt.Sprintf("%s, ваш %d на Авито", name, year),
		"subtitle": "Год находок, сделок и удачных объявлений",
		"actions": []map[string]any{
			{"type": "next", "label": "Смотреть итоги"},
		},
	})

	if m.ListingsPublished > 0 {
		story = append(story, map[string]any{
			"id":    "stat-listings",
			"type":  "stat",
			"value": "listingsPublished",
			"unit": map[string]string{
				"one": "объявление", "few": "объявления", "many": "объявлений",
			},
			"title":   "вы опубликовали",
			"eyebrow": fmt.Sprintf("За %d год", year),
		})
	}

	if m.ViewsTotal > 0 {
		story = append(story, map[string]any{
			"id":    "stat-views",
			"type":  "stat",
			"value": "viewsTotal",
			"unit": map[string]string{
				"one": "просмотр", "few": "просмотра", "many": "просмотров",
			},
			"title": "собрали ваши объявления",
		})
	}

	if m.FavoritesReceived > 0 {
		story = append(story, map[string]any{
			"id":    "stat-favorites",
			"type":  "stat",
			"value": "favoritesReceived",
			"unit": map[string]string{
				"one": "добавление", "few": "добавления", "many": "добавлений",
			},
			"title":   "в избранное",
			"eyebrow": "Любимчики покупателей",
		})
	}

	if m.MessagesSent > 0 {
		story = append(story, map[string]any{
			"id":    "stat-messages",
			"type":  "stat",
			"value": "messagesSent",
			"unit": map[string]string{
				"one": "сообщение", "few": "сообщения", "many": "сообщений",
			},
			"title":   "в чатах с покупателями",
			"eyebrow": "Диалоги",
		})
	}

	if m.DealsClosed > 0 {
		story = append(story, map[string]any{
			"id":    "stat-deals",
			"type":  "stat",
			"value": "dealsClosed",
			"unit": map[string]string{
				"one": "сделка", "few": "сделки", "many": "сделок",
			},
			"title":   "успешно закрыто",
			"eyebrow": "Результат",
		})
	}

	if m.MoneyEarned > 0 {
		story = append(story, map[string]any{
			"id":      "stat-earned",
			"type":    "stat",
			"value":   "moneyEarned",
			"unit":    "₽",
			"title":   "заработали на продажах",
			"eyebrow": "Доход",
			"valueFormat": map[string]any{
				"maximumFractionDigits": 0,
			},
		})
	}

	for _, b := range buildBadges(m) {
		story = append(story, map[string]any{
			"id":      "achievement-" + b.ID,
			"type":    "achievement",
			"badgeId": b.ID,
		})
	}

	story = append(story, map[string]any{
		"id":       "outro",
		"type":     "outro",
		"title":    "Это был ваш год на Авито",
		"subtitle": "Сохраните итоги или вернитесь к объявлениям",
		"actions": []map[string]any{
			{
				"type":    "custom",
				"id":      "close-recap",
				"label":   "На главную",
				"variant": "primary",
			},
		},
	})

	return story
}

func buildMetrics(m clients.Metrics) map[string]models.MetricValue {
	return map[string]models.MetricValue{
		"listingsPublished": {
			Type:  "number",
			Value: m.ListingsPublished,
		},
		"viewsTotal": {
			Type:  "number",
			Value: m.ViewsTotal,
		},
		"favoritesReceived": {
			Type:  "number",
			Value: m.FavoritesReceived,
		},
		"messagesSent": {
			Type:  "number",
			Value: m.MessagesSent,
		},
		"dealsClosed": {
			Type:  "number",
			Value: m.DealsClosed,
		},
		"moneyEarned": {
			Type:     "money",
			Value:    m.MoneyEarned,
			Currency: "RUB",
		},
		"moneySaved": {
			Type:     "money",
			Value:    m.MoneySaved,
			Currency: "RUB",
		},
		"daysActive": {
			Type:  "number",
			Value: m.DaysActive,
		},
		"peakDayViews": {
			Type:  "number",
			Value: m.PeakDayViews,
		},
		"searchQueries": {
			Type:  "number",
			Value: m.SearchQueries,
		},
		"categoriesTried": {
			Type:  "number",
			Value: m.CategoriesTried,
		},
		"deliveryOrders": {
			Type:  "number",
			Value: m.DeliveryOrders,
		},
	}
}

func buildBadges(metrics clients.Metrics) []models.Badge {
	var badges []models.Badge

	if metrics.MessagesSent > 1000 {
		badges = append(badges, models.Badge{
			ID:          "messages_sent_1000",
			Title:       "Разговорчивый",
			Description: fmt.Sprintf("Вы отправили %d сообщений!", metrics.MessagesSent),
		})
	}

	if metrics != (clients.Metrics{}) {
		badges = append(badges, models.Badge{
			ID:          "active_user",
			Title:       "Активный пользователь",
			Description: "Вы заходили на нашу площадку в этом году!",
		})
	}

	return badges
}
