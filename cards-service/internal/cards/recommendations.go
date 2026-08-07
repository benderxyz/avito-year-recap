package cards

import (
	"fmt"
	"sort"
	"strings"

	"cards-service/internal/clients"
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
	cond      func(clients.Metrics) bool
}

var recommendationRules = []recommendationRule{
	{
		id:        "resume-listings",
		title:     "Продвиньте свои объявления",
		text:      "У вас есть объявления, но сделок пока немного. Продвижение поможет им найти покупателя быстрее.",
		callout:   "Больше показов — больше сделок",
		linkLabel: "Продвинуть объявления",
		path:      "/promotion",
		priority:  100,
		cond: func(m clients.Metrics) bool {
			return m.ListingsPublished > 0 && m.DealsClosed == 0
		},
	},
	{
		id:        "return-favorites",
		title:     "Ваше избранное ждёт",
		text:      "Вы добавляли объявления в избранное. Загляните — часть из них ещё в продаже.",
		callout:   "Вернитесь к тому, что вам понравилось",
		linkLabel: "Открыть избранное",
		path:      "/favorites",
		priority:  90,
		cond: func(m clients.Metrics) bool {
			return m.FavoritesReceived > 0
		},
	},
	{
		id:        "resume-delivery",
		title:     "Авито Доставка",
		text:      "Вы уже пользовались доставкой. С ней проще продавать и покупать по всей стране.",
		callout:   "Продавайте покупателям из других городов",
		linkLabel: "Подключить доставку",
		path:      "/delivery",
		priority:  80,
		cond: func(m clients.Metrics) bool {
			return m.DeliveryOrders > 0
		},
	},
	{
		id:        "resume-active",
		title:     "Обновите активные объявления",
		text:      "У вас есть активные объявления. Освежите их, чтобы они снова поднялись в поиске.",
		callout:   "Актуальные объявления смотрят чаще",
		linkLabel: "К моим объявлениям",
		path:      "/profile/items",
		priority:  70,
		cond: func(m clients.Metrics) bool {
			return m.ActiveListings > 0
		},
	},
	{
		id:        "resume-search",
		title:     "Продолжите поиск",
		text:      "Вы активно искали в этом году. Возможно, нужное объявление уже появилось.",
		callout:   "Новые объявления каждый день",
		linkLabel: "Продолжить поиск",
		path:      "/",
		priority:  60,
		cond: func(m clients.Metrics) bool {
			return m.SearchQueries > 0 || m.ViewsTotal > 0
		},
	},
	{
		id:        "new-listing",
		title:     "Разместите первое объявление",
		text:      "В этом году вы ещё не продавали. Начните — это займёт пару минут.",
		callout:   "Продайте то, чем больше не пользуетесь",
		linkLabel: "Разместить объявление",
		path:      "/add",
		priority:  50,
		cond: func(m clients.Metrics) bool {
			return m.ListingsPublished == 0
		},
	},
}

func buildRecommendations(m clients.Metrics, productBaseURL string) []map[string]any {
	matched := make([]recommendationRule, 0, len(recommendationRules))
	for _, rule := range recommendationRules {
		if rule.cond(m) {
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
