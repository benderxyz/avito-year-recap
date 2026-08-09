package cards

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
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
	ruleSet := RuleSet{}
	if opts.Rules != nil {
		ruleSet = *opts.Rules
	}

	badges := buildBadges(ruleSet.badges, metrics, opts.Mode)
	slog.Debug("badges generated", "badges", badges)

	shareURL := ""
	if opts.Mode == models.RecapModePrivate && len(opts.SigningKey) > 0 {
		token := GenerateShareToken(
			opts.SigningKey,
			strconv.FormatUint(profile.UserID, 10),
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
		Metrics: buildMetrics(metrics, opts.Mode, ruleSet.metrics),
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
	defsByKey := metricDefinitionsByKey(rules.metrics)

	name := profile.Username
	if name == "" {
		name = "вы"
	}

	story := make([]map[string]any, 0)
	var outro map[string]any

	for _, rule := range rules.stories {
		if !rule.visibility.allows(mode) || !rule.when.eval(snapshot) {
			continue
		}

		scene := cloneScene(rule.scene)
		sceneType, _ := scene["type"].(string)
		metricKey := ruleMetricKey(rule)

		scene = applySceneTemplates(scene, name, year, m, metricKey)
		scene = attachPercentile(scene, m, defsByKey)

		if sceneType == "outro" {
			outro = finalizeOutro(scene, mode, shareURL)
			continue
		}

		story = append(story, scene)
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

	if outro != nil {
		story = append(story, outro)
	}

	return story
}

func ruleMetricKey(rule storyRule) string {
	if len(rule.when.predicates) == 0 {
		return ""
	}
	return string(rule.when.predicates[0].metric)
}

func applySceneTemplates(scene map[string]any, displayName string, year int, m clients.Metrics, metricKey string) map[string]any {
	replacements := map[string]string{
		"{{displayName}}": displayName,
		"{{year}}":        strconv.Itoa(year),
	}

	valueKey, _ := scene["value"].(string)
	if valueKey == "" {
		valueKey = metricKey
	}
	if valueKey != "" {
		if formatted, ok := formatMetricTemplateValue(m, valueKey); ok {
			replacements["{{value}}"] = formatted
		}
	}

	for key, raw := range scene {
		if text, ok := raw.(string); ok {
			scene[key] = applyReplacements(text, replacements)
		}
	}

	return scene
}

func formatMetricTemplateValue(m clients.Metrics, key string) (string, bool) {
	value, ok := extractMetricValue(m, key)
	if !ok {
		return "", false
	}

	switch v := value.(type) {
	case string:
		return v, true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case int:
		return strconv.Itoa(v), true
	case int64:
		return strconv.FormatInt(v, 10), true
	default:
		return fmt.Sprint(v), true
	}
}

func applyReplacements(text string, replacements map[string]string) string {
	result := text
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

func finalizeOutro(scene map[string]any, mode models.RecapMode, shareURL string) map[string]any {
	actions := normalizeActions(scene["actions"])
	scene["actions"] = filterOutroActions(actions, mode, shareURL)
	return scene
}

func normalizeActions(raw any) []map[string]any {
	switch actions := raw.(type) {
	case []map[string]any:
		return actions
	case []any:
		result := make([]map[string]any, 0, len(actions))
		for _, item := range actions {
			if action, ok := item.(map[string]any); ok {
				result = append(result, action)
			}
		}
		return result
	default:
		return nil
	}
}

func filterOutroActions(actions []map[string]any, mode models.RecapMode, shareURL string) []map[string]any {
	filtered := make([]map[string]any, 0, len(actions))
	for _, action := range actions {
		cloned := cloneScene(action)
		actionType, _ := cloned["type"].(string)
		if actionType == "share" {
			if mode != models.RecapModePrivate {
				continue
			}
			if share, ok := cloned["share"].(map[string]any); ok {
				shareClone := cloneScene(share)
				if shareURL != "" {
					shareClone["url"] = shareURL
				}
				cloned["share"] = shareClone
			}
		}
		filtered = append(filtered, cloned)
	}
	return filtered
}

func cloneScene(scene map[string]any) map[string]any {
	cloned := make(map[string]any, len(scene))
	for key, value := range scene {
		cloned[key] = value
	}
	return cloned
}

func attachPercentile(scene map[string]any, m clients.Metrics, defsByKey map[string]models.MetricDefinition) map[string]any {
	valueKey, _ := scene["value"].(string)
	if valueKey == "" {
		return scene
	}

	def, ok := defsByKey[valueKey]
	if !ok || def.PercentileKey == "" {
		return scene
	}

	if _, ok := extractMetricValue(m, def.PercentileKey); !ok {
		return scene
	}

	scene["percentile"] = def.PercentileKey
	if _, hasTemplate := scene["comparisonTemplate"]; !hasTemplate {
		scene["comparisonTemplate"] = "это больше, чем у {{percentile}}% пользователей"
	}
	return scene
}

func metricDefinitionsByKey(defs []models.MetricDefinition) map[string]models.MetricDefinition {
	result := make(map[string]models.MetricDefinition, len(defs))
	for _, def := range defs {
		result[def.Key] = def
	}
	return result
}

func buildMetrics(m clients.Metrics, mode models.RecapMode, defs []models.MetricDefinition) map[string]models.MetricValue {
	all := make(map[string]models.MetricValue, len(defs))

	for _, def := range defs {
		value, ok := extractMetricValue(m, def.Key)
		if !ok {
			continue
		}

		metric := models.MetricValue{
			Type:  def.ValueType,
			Value: value,
		}
		if def.ValueType == models.MetricTypeMoney && def.Currency != "" {
			metric.Currency = def.Currency
		}
		all[def.Key] = metric
	}

	if mode == models.RecapModePublic {
		filtered := make(map[string]models.MetricValue, len(all))
		for _, def := range defs {
			if !def.IsPublic {
				continue
			}
			if value, ok := all[def.Key]; ok {
				filtered[def.Key] = value
			}
		}
		return filtered
	}

	return all
}

func extractMetricValue(m clients.Metrics, key string) (any, bool) {
	switch models.MetricKey(key) {
	case models.MetricListingsPublished:
		return m.ListingsPublished, true
	case models.MetricListingsPercentile:
		if m.ListingsPercentile == nil {
			return nil, false
		}
		return *m.ListingsPercentile, true
	case models.MetricViewsTotal:
		return m.ViewsTotal, true
	case models.MetricViewsPercentile:
		if m.ViewsPercentile == nil {
			return nil, false
		}
		return *m.ViewsPercentile, true
	case models.MetricFavoritesReceived:
		return m.FavoritesReceived, true
	case models.MetricFavoritesPercentile:
		if m.FavoritesPercentile == nil {
			return nil, false
		}
		return *m.FavoritesPercentile, true
	case models.MetricMessagesSent:
		return m.MessagesSent, true
	case models.MetricMessagesPercentile:
		if m.MessagesPercentile == nil {
			return nil, false
		}
		return *m.MessagesPercentile, true
	case models.MetricDealsClosed:
		return m.DealsClosed, true
	case models.MetricDealsPercentile:
		if m.DealsPercentile == nil {
			return nil, false
		}
		return *m.DealsPercentile, true
	case models.MetricMoneyEarned:
		return m.MoneyEarned, true
	case models.MetricMoneySaved:
		return m.MoneySaved, true
	case models.MetricDaysActive:
		return m.DaysActive, true
	case models.MetricPeakDayViews:
		return m.PeakDayViews, true
	case models.MetricSearchQueries:
		return m.SearchQueries, true
	case models.MetricCategoriesTried:
		return m.CategoriesTried, true
	case models.MetricDeliveryOrders:
		return m.DeliveryOrders, true
	case models.MetricActiveListings:
		return m.ActiveListings, true
	case models.MetricSellerRating:
		return m.SellerRating, true
	case models.MetricAvgReplySeconds:
		return m.AvgReplySeconds, true
	case models.MetricFirstListingAt:
		if m.FirstListingAt == 0 {
			return nil, false
		}
		return time.Unix(m.FirstListingAt, 0).UTC().Format("2006-01-02"), true
	case models.MetricFirstDealAt:
		if m.FirstDealAt == 0 {
			return nil, false
		}
		return time.Unix(m.FirstDealAt, 0).UTC().Format("2006-01-02"), true
	default:
		return nil, false
	}
}
