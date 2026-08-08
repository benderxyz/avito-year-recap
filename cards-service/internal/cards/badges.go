package cards

import (
	"cards-service/internal/clients"
	"cards-service/internal/models"
)

type badgeRule struct {
	badge      models.Badge
	visibility visibility
	when       condition
}

var defaultBadgeRules = []badgeRule{
	{
		badge: models.Badge{
			ID:          "messages_sent_1000",
			Title:       "Разговорчивый",
			Description: "Вы отправили больше 1000 сообщений!",
			IconURL:     "/badges/messages_sent_1000.svg",
		},
		visibility: visibilityBoth,
		when: condition{predicates: []predicate{
			{metric: models.MetricMessagesSent, op: opGt, value: 1000},
		}},
	},
	{
		badge: models.Badge{
			ID:          "active_user",
			Title:       "Активный пользователь",
			Description: "Вы заходили на нашу площадку в этом году!",
			IconURL:     "/badges/active_user.svg",
		},
		visibility: visibilityBoth,
		when: condition{predicates: []predicate{
			{metric: models.MetricDaysActive, op: opGt, value: 0},
		}},
	},
}

func buildBadges(rules []badgeRule, m clients.Metrics, mode models.RecapMode) []models.Badge {
	snapshot := metricsSnapshot(m)

	var badges []models.Badge
	for _, rule := range rules {
		if !rule.visibility.allows(mode) {
			continue
		}
		if rule.when.eval(snapshot) {
			badges = append(badges, rule.badge)
		}
	}

	return badges
}
