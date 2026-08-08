package cards

import (
	"fmt"
	"time"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

type BuildOptions struct {
	Mode           models.RecapMode
	SigningKey     []byte
	ShareBaseURL   string
	ProductBaseURL string
	Rules          *RuleSet
}

func BuildRecap(
	profile models.Profile,
	year int,
	metrics clients.Metrics,
	opts BuildOptions,
) models.RecapPayload {
	ruleSet := resolveRuleSet(opts.Rules)

	badges := buildBadges(ruleSet.badges, metrics, opts.Mode)

	shareURL := ""
	if opts.Mode == models.RecapModePrivate && len(opts.SigningKey) > 0 {
		token := GenerateShareToken(
			opts.SigningKey,
			profile.ExternalID,
			year,
		)
		shareURL = buildShareURL(opts.ShareBaseURL, token)
	}

	payload := models.RecapPayload{
		SchemaVersion: 1,
		Meta: models.Meta{
			Vertical: "marketplace",
			Year:     year,
			Locale:   "ru-RU",
			User: models.User{
				ID:          profile.ExternalID,
				DisplayName: profile.Username,
			},
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		},
		Metrics: buildMetrics(metrics, opts.Mode),
		Badges:  badges,
		Story:   buildStory(profile, year, metrics, ruleSet, badges, opts, shareURL),
	}

	if shareURL != "" {
		payload.Features = &models.Features{
			ShareEnabled: true,
			ShareURL:     shareURL,
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

func buildStory(profile models.Profile, year int, m clients.Metrics, rules RuleSet, badges []models.Badge, opts BuildOptions, shareURL string) []map[string]any {
	mode := opts.Mode
	snapshot := metricsSnapshot(m)

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

	for _, rule := range rules.stories {
		if rule.visibility.allows(mode) && rule.when.eval(snapshot) {
			story = append(story, attachPercentile(cloneScene(rule.scene), m))
		}
	}

	for _, b := range badges {
		story = append(story, map[string]any{
			"id":      "achievement-" + b.ID,
			"type":    "achievement",
			"badgeId": b.ID,
		})
	}

	if mode == models.RecapModePrivate {
		story = append(story, buildRecommendations(rules.recommendations, m, opts.ProductBaseURL)...)
	}

	story = append(story, buildOutro(mode, shareURL))

	return story
}

func buildOutro(mode models.RecapMode, shareURL string) map[string]any {
	actions := make([]map[string]any, 0, 2)

	if mode == models.RecapModePrivate {
		share := map[string]any{
			"kind":  "link",
			"title": "Мои итоги на Авито",
			"text":  "Посмотрите, каким был мой год на Авито!",
		}
		if shareURL != "" {
			share["url"] = shareURL
		}

		actions = append(actions, map[string]any{
			"type":  "share",
			"label": "Поделиться",
			"share": share,
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

func cloneScene(scene map[string]any) map[string]any {
	cloned := make(map[string]any, len(scene))
	for key, value := range scene {
		cloned[key] = value
	}
	return cloned
}

func attachPercentile(scene map[string]any, m clients.Metrics) map[string]any {
	valueKey, _ := scene["value"].(string)

	var percentileKey models.MetricKey
	var percentile *float64

	switch valueKey {
	case string(models.MetricListingsPublished):
		percentileKey = models.MetricListingsPercentile
		percentile = m.ListingsPercentile
	case string(models.MetricViewsTotal):
		percentileKey = models.MetricViewsPercentile
		percentile = m.ViewsPercentile
	case string(models.MetricFavoritesReceived):
		percentileKey = models.MetricFavoritesPercentile
		percentile = m.FavoritesPercentile
	case string(models.MetricMessagesSent):
		percentileKey = models.MetricMessagesPercentile
		percentile = m.MessagesPercentile
	case string(models.MetricDealsClosed):
		percentileKey = models.MetricDealsPercentile
		percentile = m.DealsPercentile
	default:
		return scene
	}

	if percentile != nil {
		scene["percentile"] = string(percentileKey)
		scene["comparisonTemplate"] = "это больше, чем у {{percentile}}% пользователей"
	}
	return scene
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

	addPercentileMetric(all, models.MetricListingsPercentile, m.ListingsPercentile)
	addPercentileMetric(all, models.MetricViewsPercentile, m.ViewsPercentile)
	addPercentileMetric(all, models.MetricFavoritesPercentile, m.FavoritesPercentile)
	addPercentileMetric(all, models.MetricMessagesPercentile, m.MessagesPercentile)
	addPercentileMetric(all, models.MetricDealsPercentile, m.DealsPercentile)

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

func addPercentileMetric(
	all map[string]models.MetricValue,
	key models.MetricKey,
	value *float64,
) {
	if value == nil {
		return
	}
	all[string(key)] = models.MetricValue{
		Type:  models.MetricTypePercentile,
		Value: *value,
	}
}
