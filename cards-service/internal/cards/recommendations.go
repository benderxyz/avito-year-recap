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
	when      condition
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
