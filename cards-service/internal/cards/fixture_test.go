package cards

import (
	"cards-service/internal/models"
)

func testMetricDefinitions() []models.MetricDefinition {
	number := func(key string, isPublic bool, percentileKey string) models.MetricDefinition {
		return models.MetricDefinition{
			Key:                     key,
			ValueType:               models.MetricTypeNumber,
			IsPublic:                isPublic,
			PercentileKey:           percentileKey,
			ComparisonMinPercentile: defaultComparisonMinPercentile,
			SourceKey:               key,
			SourceField:             models.MetricSourceValue,
			IncludeInLLM:            true,
		}
	}

	percentile := func(key, sourceKey string, isPublic bool) models.MetricDefinition {
		return models.MetricDefinition{
			Key:         key,
			ValueType:   models.MetricTypePercentile,
			IsPublic:    isPublic,
			SourceKey:   sourceKey,
			SourceField: models.MetricSourcePercentile,
		}
	}

	money := func(key string) models.MetricDefinition {
		return models.MetricDefinition{
			Key:          key,
			ValueType:    models.MetricTypeMoney,
			Currency:     models.CurrencyRUB,
			SourceKey:    key,
			SourceField:  models.MetricSourceValue,
			IncludeInLLM: true,
		}
	}

	date := func(key string) models.MetricDefinition {
		return models.MetricDefinition{
			Key:         key,
			ValueType:   models.MetricTypeDate,
			SourceKey:   key,
			SourceField: models.MetricSourceValue,
		}
	}

	return []models.MetricDefinition{
		number("listingsPublished", true, "listingsPercentile"),
		percentile("listingsPercentile", "listingsPublished", true),
		number("viewsTotal", true, "viewsPercentile"),
		percentile("viewsPercentile", "viewsTotal", true),
		number("favoritesReceived", true, "favoritesPercentile"),
		percentile("favoritesPercentile", "favoritesReceived", true),
		number("messagesSent", false, "messagesPercentile"),
		percentile("messagesPercentile", "messagesSent", false),
		number("dealsClosed", true, "dealsPercentile"),
		percentile("dealsPercentile", "dealsClosed", true),
		money("moneyEarned"),
		money("moneySaved"),
		number("daysActive", true, ""),
		number("peakDayViews", true, ""),
		number("categoriesTried", true, ""),
		number("searchQueries", true, ""),
		number("deliveryOrders", true, ""),
		number("activeListings", false, ""),
		number("avgReplySeconds", false, ""),
		date("firstListingAt"),
		date("firstDealAt"),
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
		must("both", "listingsPublished", "gt", 0, true, `{"id":"stat-listings","type":"stat","value":"listingsPublished","unit":{"one":"объявление","few":"объявления","many":"объявлений"},"title":"вы опубликовали","eyebrow":"За год","comparisonTemplate":"Объявлений у вас больше, чем у {{percentile}}% продавцов"}`),
		must("both", "viewsTotal", "gt", 0, true, `{"id":"stat-views","type":"stat","value":"viewsTotal","title":"собрали ваши объявления","comparisonTemplate":"Ваши объявления смотрели чаще, чем у {{percentile}}% продавцов"}`),
		must("both", "favoritesReceived", "gt", 0, true, `{"id":"stat-favorites","type":"stat","value":"favoritesReceived","title":"в избранное","comparisonTemplate":"В избранное вас добавляли чаще, чем у {{percentile}}% продавцов"}`),
		must("private", "messagesSent", "gt", 0, true, `{"id":"stat-messages","type":"stat","value":"messagesSent","title":"в чатах с покупателями","comparisonTemplate":"Вы переписывались активнее, чем {{percentile}}% пользователей"}`),
		must("both", "dealsClosed", "gt", 0, true, `{"id":"stat-deals","type":"stat","value":"dealsClosed","title":"успешно закрыто","comparisonTemplate":"Сделок больше, чем у {{percentile}}% продавцов"}`),
		must("private", "moneyEarned", "gt", 0, true, `{"id":"stat-earned","type":"stat","value":"moneyEarned","title":"заработали на продажах"}`),
		must("private", "moneySaved", "gt", 0, true, `{"id":"stat-saved","type":"stat","value":"moneySaved","title":"сэкономили на покупках"}`),
		must("both", "daysActive", "gt", 0, true, `{"id":"blocks-days-active","type":"blocks","blocks":[{"type":"stat","value":"daysActive","title":"были активны"}]}`),
		must("both", "peakDayViews", "gt", 0, true, `{"id":"stat-peak-views","type":"stat","value":"peakDayViews","title":"за один день"}`),
		must("both", "categoriesTried", "gt", 0, true, `{"id":"stat-categories","type":"blocks","blocks":[{"type":"stat","value":"categoriesTried","unit":{"one":"категория","few":"категории","many":"категорий"},"title":"вы открывали за год","eyebrow":"Кругозор"},{"type":"callout","text":"Заглядывали в разные разделы, от электроники до садовой мебели"}]}`),
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
