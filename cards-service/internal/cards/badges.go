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
