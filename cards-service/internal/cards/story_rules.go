package cards

import "cards-service/internal/models"

type storyRule struct {
	visibility visibility
	when       condition
	scene      map[string]any
}

var defaultStatSceneRules = []storyRule{
	{
		visibility: visibilityBoth,
		when: condition{predicates: []predicate{
			{metric: models.MetricListingsPublished, op: opGt, value: 0},
		}},
		scene: map[string]any{
			"id":    "stat-listings",
			"type":  "stat",
			"value": "listingsPublished",
			"unit": map[string]string{
				"one": "объявление", "few": "объявления", "many": "объявлений",
			},
			"title":   "вы опубликовали",
			"eyebrow": "За год",
		},
	},
	{
		visibility: visibilityBoth,
		when: condition{predicates: []predicate{
			{metric: models.MetricViewsTotal, op: opGt, value: 0},
		}},
		scene: map[string]any{
			"id":    "stat-views",
			"type":  "stat",
			"value": "viewsTotal",
			"unit": map[string]string{
				"one": "просмотр", "few": "просмотра", "many": "просмотров",
			},
			"title": "собрали ваши объявления",
		},
	},
	{
		visibility: visibilityBoth,
		when: condition{predicates: []predicate{
			{metric: models.MetricFavoritesReceived, op: opGt, value: 0},
		}},
		scene: map[string]any{
			"id":    "stat-favorites",
			"type":  "stat",
			"value": "favoritesReceived",
			"unit": map[string]string{
				"one": "добавление", "few": "добавления", "many": "добавлений",
			},
			"title":   "в избранное",
			"eyebrow": "Любимчики покупателей",
		},
	},
	{
		visibility: visibilityPrivate,
		when: condition{predicates: []predicate{
			{metric: models.MetricMessagesSent, op: opGt, value: 0},
		}},
		scene: map[string]any{
			"id":    "stat-messages",
			"type":  "stat",
			"value": "messagesSent",
			"unit": map[string]string{
				"one": "сообщение", "few": "сообщения", "many": "сообщений",
			},
			"title":   "в чатах с покупателями",
			"eyebrow": "Диалоги",
		},
	},
	{
		visibility: visibilityBoth,
		when: condition{predicates: []predicate{
			{metric: models.MetricDealsClosed, op: opGt, value: 0},
		}},
		scene: map[string]any{
			"id":    "stat-deals",
			"type":  "stat",
			"value": "dealsClosed",
			"unit": map[string]string{
				"one": "сделка", "few": "сделки", "many": "сделок",
			},
			"title":   "успешно закрыто",
			"eyebrow": "Результат",
		},
	},
	{
		visibility: visibilityPrivate,
		when: condition{predicates: []predicate{
			{metric: models.MetricMoneyEarned, op: opGt, value: 0},
		}},
		scene: map[string]any{
			"id":      "stat-earned",
			"type":    "stat",
			"value":   "moneyEarned",
			"unit":    "₽",
			"title":   "заработали на продажах",
			"eyebrow": "Доход",
			"valueFormat": map[string]any{
				"maximumFractionDigits": 0,
			},
		},
	},
}
