package cards

import (
	"fmt"
	"sort"
	"strings"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

const maxRecommendations = 3

type recommendationRule struct {
	id        string
	title     string
	text      string
	callout   string
	linkLabel string
	path      string
	priority  int
	when      condition
}

var defaultRecommendationRules = []recommendationRule{
	{
		id:        "resume-listings",
		title:     "Продвиньте свои объявления",
		text:      "У вас есть объявления, но сделок пока немного. Продвижение поможет им найти покупателя быстрее.",
		callout:   "Больше показов — больше сделок",
		linkLabel: "Продвинуть объявления",
		path:      "/profile",
		priority:  100,
		when: condition{match: matchAll, predicates: []predicate{
			{metric: models.MetricListingsPublished, op: opGt, value: 0},
			{metric: models.MetricDealsClosed, op: opEq, value: 0},
		}},
	},
	{
		id:        "listings-favorited",
		title:     "Ваши объявления в избранном",
		text:      "Покупатели добавляли ваши объявления в избранное. Освежите их и ответьте на вопросы, пока интерес высокий.",
		callout:   "Вами интересуются покупатели",
		linkLabel: "К моим объявлениям",
		path:      "/profile",
		priority:  90,
		when: condition{predicates: []predicate{
			{metric: models.MetricFavoritesReceived, op: opGt, value: 0},
		}},
	},
	{
		id:        "resume-delivery",
		title:     "Авито Доставка",
		text:      "Вы уже пользовались доставкой. С ней проще продавать и покупать по всей стране.",
		callout:   "Продавайте покупателям из других городов",
		linkLabel: "Подключить доставку",
		path:      "/dostavka",
		priority:  80,
		when: condition{predicates: []predicate{
			{metric: models.MetricDeliveryOrders, op: opGt, value: 0},
		}},
	},
	{
		id:        "resume-active",
		title:     "Обновите активные объявления",
		text:      "У вас есть активные объявления. Освежите их, чтобы они снова поднялись в поиске.",
		callout:   "Актуальные объявления смотрят чаще",
		linkLabel: "К моим объявлениям",
		path:      "/profile",
		priority:  70,
		when: condition{predicates: []predicate{
			{metric: models.MetricActiveListings, op: opGt, value: 0},
		}},
	},
	{
		id:        "resume-search",
		title:     "Продолжите поиск",
		text:      "Вы активно искали в этом году. Возможно, нужное объявление уже появилось.",
		callout:   "Новые объявления каждый день",
		linkLabel: "Продолжить поиск",
		path:      "/",
		priority:  60,
		when: condition{match: matchAny, predicates: []predicate{
			{metric: models.MetricSearchQueries, op: opGt, value: 0},
			{metric: models.MetricViewsTotal, op: opGt, value: 0},
		}},
	},
	{
		id:        "new-listing",
		title:     "Разместите первое объявление",
		text:      "В этом году вы ещё не продавали. Начните — это займёт пару минут.",
		callout:   "Продайте то, чем больше не пользуетесь",
		linkLabel: "Разместить объявление",
		path:      "/additem",
		priority:  50,
		when: condition{predicates: []predicate{
			{metric: models.MetricListingsPublished, op: opEq, value: 0},
		}},
	},
}

func buildRecommendations(rules []recommendationRule, m clients.Metrics, productBaseURL string) []map[string]any {
	snapshot := metricsSnapshot(m)

	matched := make([]recommendationRule, 0, len(rules))
	for _, rule := range rules {
		if rule.when.eval(snapshot) {
			matched = append(matched, rule)
		}
	}

	sort.SliceStable(matched, func(i, j int) bool {
		return matched[i].priority > matched[j].priority
	})

	if len(matched) > maxRecommendations {
		matched = matched[:maxRecommendations]
	}

	scenes := make([]map[string]any, 0, len(matched))
	for _, rule := range matched {
		scenes = append(scenes, map[string]any{
			"id":      "reco-" + rule.id,
			"type":    "upsell",
			"title":   rule.title,
			"text":    rule.text,
			"callout": rule.callout,
			"actions": []map[string]any{
				{
					"type":    "link",
					"label":   rule.linkLabel,
					"href":    buildProductURL(productBaseURL, rule.path),
					"target":  "_blank",
					"variant": "primary",
				},
				{
					"type":    "next",
					"label":   "Пропустить",
					"variant": "ghost",
				},
			},
		})
	}

	return scenes
}

func buildProductURL(baseURL, path string) string {
	if baseURL == "" {
		return path
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(path, "/"))
}
