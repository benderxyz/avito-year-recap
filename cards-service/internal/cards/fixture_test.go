package cards

import (
	"cards-service/internal/models"
)

func testMetricDefinitions() []models.MetricDefinition {
	return []models.MetricDefinition{
		{Key: "listingsPublished", ValueType: models.MetricTypeNumber, IsPublic: true, PercentileKey: "listingsPercentile"},
		{Key: "listingsPercentile", ValueType: models.MetricTypePercentile, IsPublic: true},
		{Key: "viewsTotal", ValueType: models.MetricTypeNumber, IsPublic: true, PercentileKey: "viewsPercentile"},
		{Key: "viewsPercentile", ValueType: models.MetricTypePercentile, IsPublic: true},
		{Key: "favoritesReceived", ValueType: models.MetricTypeNumber, IsPublic: true, PercentileKey: "favoritesPercentile"},
		{Key: "favoritesPercentile", ValueType: models.MetricTypePercentile, IsPublic: true},
		{Key: "messagesSent", ValueType: models.MetricTypeNumber, IsPublic: false, PercentileKey: "messagesPercentile"},
		{Key: "messagesPercentile", ValueType: models.MetricTypePercentile, IsPublic: false},
		{Key: "dealsClosed", ValueType: models.MetricTypeNumber, IsPublic: true, PercentileKey: "dealsPercentile"},
		{Key: "dealsPercentile", ValueType: models.MetricTypePercentile, IsPublic: true},
		{Key: "moneyEarned", ValueType: models.MetricTypeMoney, Currency: models.CurrencyRUB, IsPublic: false},
		{Key: "moneySaved", ValueType: models.MetricTypeMoney, Currency: models.CurrencyRUB, IsPublic: false},
		{Key: "daysActive", ValueType: models.MetricTypeNumber, IsPublic: true},
		{Key: "peakDayViews", ValueType: models.MetricTypeNumber, IsPublic: true},
		{Key: "categoriesTried", ValueType: models.MetricTypeNumber, IsPublic: true},
		{Key: "searchQueries", ValueType: models.MetricTypeNumber, IsPublic: true},
		{Key: "deliveryOrders", ValueType: models.MetricTypeNumber, IsPublic: true},
		{Key: "activeListings", ValueType: models.MetricTypeNumber, IsPublic: false},
		{Key: "sellerRating", ValueType: models.MetricTypeNumber, IsPublic: false},
		{Key: "avgReplySeconds", ValueType: models.MetricTypeNumber, IsPublic: false},
		{Key: "firstListingAt", ValueType: models.MetricTypeString, IsPublic: false},
		{Key: "firstDealAt", ValueType: models.MetricTypeString, IsPublic: false},
	}
}

func testStoryRules() []storyRule {
	must := func(vis, metricKey, op string, threshold float64, hasPredicate bool, payload string) storyRule {
		rule, err := makeStoryRule(vis, metricKey, op, threshold, hasPredicate, []byte(payload))
		if err != nil {
			panic(err)
		}
		return rule
	}

	return []storyRule{
		must("both", "", "", 0, false, `{"id":"intro","type":"intro","title":"{{displayName}}, ваш {{year}} на Авито","subtitle":"Год находок, сделок и удачных объявлений","actions":[{"type":"next","label":"Смотреть итоги"}]}`),
		must("both", "listingsPublished", "gt", 0, true, `{"id":"stat-listings","type":"stat","value":"listingsPublished","unit":{"one":"объявление","few":"объявления","many":"объявлений"},"title":"вы опубликовали","eyebrow":"За год"}`),
		must("both", "viewsTotal", "gt", 0, true, `{"id":"stat-views","type":"stat","value":"viewsTotal","title":"собрали ваши объявления"}`),
		must("both", "favoritesReceived", "gt", 0, true, `{"id":"stat-favorites","type":"stat","value":"favoritesReceived","title":"в избранное"}`),
		must("private", "messagesSent", "gt", 0, true, `{"id":"stat-messages","type":"stat","value":"messagesSent","title":"в чатах с покупателями"}`),
		must("both", "dealsClosed", "gt", 0, true, `{"id":"stat-deals","type":"stat","value":"dealsClosed","title":"успешно закрыто"}`),
		must("private", "moneyEarned", "gt", 0, true, `{"id":"stat-earned","type":"stat","value":"moneyEarned","title":"заработали на продажах"}`),
		must("private", "moneySaved", "gt", 0, true, `{"id":"stat-saved","type":"stat","value":"moneySaved","title":"сэкономили на покупках"}`),
		must("both", "daysActive", "gt", 0, true, `{"id":"blocks-days-active","type":"blocks","blocks":[{"type":"stat","value":"daysActive","title":"были активны"}]}`),
		must("both", "peakDayViews", "gt", 0, true, `{"id":"stat-peak-views","type":"stat","value":"peakDayViews","title":"за один день"}`),
		must("both", "categoriesTried", "gt", 0, true, `{"id":"stat-categories","type":"stat","value":"categoriesTried","title":"вы попробовали"}`),
		must("both", "searchQueries", "gt", 0, true, `{"id":"blocks-search","type":"blocks","blocks":[{"type":"stat","value":"searchQueries","title":"запросов за год"}]}`),
		must("both", "deliveryOrders", "gt", 0, true, `{"id":"stat-delivery","type":"stat","value":"deliveryOrders","title":"через Авито Доставку"}`),
		must("private", "activeListings", "gt", 0, true, `{"id":"stat-active-listings","type":"stat","value":"activeListings","title":"сейчас активны"}`),
		must("private", "avgReplySeconds", "gt", 0, true, `{"id":"stat-reply","type":"stat","value":"avgReplySeconds","title":"в среднем до ответа"}`),
		must("private", "firstListingAt", "gt", 0, true, `{"id":"insight-first-listing","type":"insight","title":"Первое объявление","text":"Первое объявление в этом году вы опубликовали {{value}}"}`),
		must("private", "firstDealAt", "gt", 0, true, `{"id":"insight-first-deal","type":"insight","title":"Первая сделка","text":"Первую сделку закрыли {{value}}"}`),
		must("both", "", "", 0, false, `{"id":"outro","type":"outro","title":"Это был ваш год на Авито","subtitle":"Сохраните итоги или вернитесь к объявлениям","actions":[{"type":"share","label":"Поделиться","share":{"kind":"link","title":"Мои итоги на Авито","text":"Посмотрите, каким был мой год на Авито!"}},{"type":"custom","id":"close-recap","label":"На главную","variant":"primary"}]}`),
	}
}

func testBadgeRules() []badgeRule {
	return []badgeRule{
		makeBadgeRule("messages_sent_1000", "Разговорчивый", "Вы отправили больше 1000 сообщений!", "/badges/messages_sent_1000.svg", "both", "messagesSent", "gt", 1000),
		makeBadgeRule("active_user", "Активный пользователь", "Вы заходили на нашу площадку в этом году!", "/badges/active_user.svg", "both", "daysActive", "gt", 0),
	}
}

func testRecommendationRules() []recommendationRule {
	must := func(id, title, text, callout, linkLabel, path string, priority int, condition string) recommendationRule {
		rule, err := makeRecommendationRule(id, title, text, callout, linkLabel, path, priority, []byte(condition))
		if err != nil {
			panic(err)
		}
		return rule
	}

	return []recommendationRule{
		must("resume-listings", "Продвиньте свои объявления", "text", "callout", "Продвинуть объявления", "/profile", 100,
			`{"match":"all","predicates":[{"metric":"listingsPublished","op":"gt","value":0},{"metric":"dealsClosed","op":"eq","value":0}]}`),
		must("listings-favorited", "Ваши объявления в избранном", "text", "callout", "К моим объявлениям", "/profile", 90,
			`{"predicates":[{"metric":"favoritesReceived","op":"gt","value":0}]}`),
		must("resume-delivery", "Авито Доставка", "text", "callout", "Подключить доставку", "/dostavka", 80,
			`{"predicates":[{"metric":"deliveryOrders","op":"gt","value":0}]}`),
		must("resume-active", "Обновите активные объявления", "text", "callout", "К моим объявлениям", "/profile", 70,
			`{"predicates":[{"metric":"activeListings","op":"gt","value":0}]}`),
		must("resume-search", "Продолжите поиск", "text", "callout", "Продолжить поиск", "/", 60,
			`{"match":"any","predicates":[{"metric":"searchQueries","op":"gt","value":0},{"metric":"viewsTotal","op":"gt","value":0}]}`),
		must("new-listing", "Разместите первое объявление", "text", "callout", "Разместить объявление", "/additem", 50,
			`{"predicates":[{"metric":"listingsPublished","op":"eq","value":0}]}`),
	}
}

func testRuleSet() RuleSet {
	return RuleSet{
		badges:          testBadgeRules(),
		stories:         testStoryRules(),
		recommendations: testRecommendationRules(),
		metrics:         testMetricDefinitions(),
	}
}
