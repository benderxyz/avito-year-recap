package cards

import (
	"fmt"
	"time"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

type BuildOptions struct {
	Mode         models.RecapMode
	SigningKey   []byte
	ShareBaseURL string
}

func BuildRecap(
	profile models.Profile,
	year int,
	metrics clients.Metrics,
	opts BuildOptions,
) models.RecapPayload {
	payload := models.RecapPayload{
		SchemaVersion: 1,
		Meta: models.Meta{
			Vertical:    "marketplace",
			Year:        year,
			Locale:      "ru-RU",
			Profile:     profile,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		},
		Metrics: buildMetrics(metrics, opts.Mode),
		Badges:  buildBadges(metrics, opts.Mode),
		Story:   buildStory(profile, year, metrics, opts.Mode),
	}

	if opts.Mode == models.RecapModePrivate && len(opts.SigningKey) > 0 {
		token := GenerateShareToken(
			opts.SigningKey,
			profile.ExternalID,
			year,
		)

		payload.Features = &models.Features{
			ShareEnabled: true,
			ShareURL:     buildShareURL(opts.ShareBaseURL, token),
		}
	}

	return payload
}

func buildShareURL(baseURL, token string) string {
	if baseURL == "" {
		return "/share/" + token
	}
	return fmt.Sprintf("%s/share/%s", baseURL, token)
}

func buildStory(profile models.Profile, year int, m clients.Metrics, mode models.RecapMode) []map[string]any {
	name := profile.Username
	if name == "" {
		name = "вы"
	}

	isPrivate := mode == models.RecapModePrivate

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

	if isPrivate && m.MessagesSent > 0 {
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

	if isPrivate && m.MoneyEarned > 0 {
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

	for _, b := range buildBadges(m, mode) {
		story = append(story, map[string]any{
			"id":      "achievement-" + b.ID,
			"type":    "achievement",
			"badgeId": b.ID,
		})
	}

	story = append(story, buildOutro(mode))

	return story
}

func buildOutro(mode models.RecapMode) map[string]any {
	actions := make([]map[string]any, 0, 2)

	if mode == models.RecapModePrivate {
		actions = append(actions, map[string]any{
			"type":  "share",
			"label": "Поделиться",
			"share": map[string]any{
				"kind":  "link",
				"title": "Мои итоги на Авито",
				"text":  "Посмотрите, каким был мой год на Авито!",
			},
		})
	}

	actions = append(actions, map[string]any{
		"type":    "custom",
		"id":      "close-recap",
		"label":   "На главную",
		"variant": "primary",
	})

	return map[string]any{
		"id":       "outro",
		"type":     "outro",
		"title":    "Это был ваш год на Авито",
		"subtitle": "Сохраните итоги или вернитесь к объявлениям",
		"actions":  actions,
	}
}

func buildMetrics(m clients.Metrics, mode models.RecapMode) map[string]models.MetricValue {
	all := map[string]models.MetricValue{
		string(models.MetricListingsPublished): {
			Type:  models.MetricTypeNumber,
			Value: m.ListingsPublished,
		},
		string(models.MetricViewsTotal): {
			Type:  models.MetricTypeNumber,
			Value: m.ViewsTotal,
		},
		string(models.MetricFavoritesReceived): {
			Type:  models.MetricTypeNumber,
			Value: m.FavoritesReceived,
		},
		string(models.MetricMessagesSent): {
			Type:  models.MetricTypeNumber,
			Value: m.MessagesSent,
		},
		string(models.MetricDealsClosed): {
			Type:  models.MetricTypeNumber,
			Value: m.DealsClosed,
		},
		string(models.MetricMoneyEarned): {
			Type:     models.MetricTypeMoney,
			Value:    m.MoneyEarned,
			Currency: models.CurrencyRUB,
		},
		string(models.MetricMoneySaved): {
			Type:     models.MetricTypeMoney,
			Value:    m.MoneySaved,
			Currency: models.CurrencyRUB,
		},
		string(models.MetricDaysActive): {
			Type:  models.MetricTypeNumber,
			Value: m.DaysActive,
		},
		string(models.MetricPeakDayViews): {
			Type:  models.MetricTypeNumber,
			Value: m.PeakDayViews,
		},
		string(models.MetricSearchQueries): {
			Type:  models.MetricTypeNumber,
			Value: m.SearchQueries,
		},
		string(models.MetricCategoriesTried): {
			Type:  models.MetricTypeNumber,
			Value: m.CategoriesTried,
		},
		string(models.MetricDeliveryOrders): {
			Type:  models.MetricTypeNumber,
			Value: m.DeliveryOrders,
		},
	}

	if mode == models.RecapModePublic {
		filtered := make(map[string]models.MetricValue, len(all))
		for key, value := range all {
			if models.PublicMetricsAllowlist[key] {
				filtered[key] = value
			}
		}
		return filtered
	}

	return all
}

func buildBadges(metrics clients.Metrics, mode models.RecapMode) []models.Badge {
	var badges []models.Badge

	if metrics.MessagesSent > 1000 {
		description := "Вы отправили больше 1000 сообщений!"
		if mode == models.RecapModePrivate {
			description = fmt.Sprintf("Вы отправили %d сообщений!", metrics.MessagesSent)
		}
		badges = append(badges, models.Badge{
			ID:          "messages_sent_1000",
			Title:       "Разговорчивый",
			Description: description,
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
