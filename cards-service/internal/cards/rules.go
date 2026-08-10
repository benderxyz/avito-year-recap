package cards

import (
	"cards-service/internal/clients"
	"cards-service/internal/models"
)

type predicateOp string

const (
	opGt     predicateOp = "gt"
	opGte    predicateOp = "gte"
	opEq     predicateOp = "eq"
	opExists predicateOp = "exists"
)

type matchMode string

const (
	matchAll matchMode = "all"
	matchAny matchMode = "any"
)

type visibility string

const (
	visibilityPrivate visibility = "private"
	visibilityPublic  visibility = "public"
	visibilityBoth    visibility = "both"
)

type predicate struct {
	metric models.MetricKey
	op     predicateOp
	value  float64
}

type condition struct {
	match      matchMode
	predicates []predicate
}

type storyRule struct {
	visibility visibility
	when       condition
	scene      map[string]any
}

type RuleSet struct {
	badges          []badgeRule
	stories         []storyRule
	recommendations []recommendationRule
	metrics         []models.MetricDefinition
}

func (r RuleSet) PublicMetricKeys() map[string]bool {
	keys := make(map[string]bool, len(r.metrics))
	for _, def := range r.metrics {
		if def.IsPublic {
			keys[def.Key] = true
		}
	}
	return keys
}

func metricsSnapshot(m clients.Metrics) map[models.MetricKey]float64 {
	return map[models.MetricKey]float64{
		models.MetricListingsPublished: float64(m.ListingsPublished),
		models.MetricViewsTotal:        float64(m.ViewsTotal),
		models.MetricFavoritesReceived: float64(m.FavoritesReceived),
		models.MetricMessagesSent:      float64(m.MessagesSent),
		models.MetricDealsClosed:       float64(m.DealsClosed),
		models.MetricMoneyEarned:       float64(m.MoneyEarned),
		models.MetricMoneySaved:        float64(m.MoneySaved),
		models.MetricDaysActive:        float64(m.DaysActive),
		models.MetricPeakDayViews:      float64(m.PeakDayViews),
		models.MetricSearchQueries:     float64(m.SearchQueries),
		models.MetricCategoriesTried:   float64(m.CategoriesTried),
		models.MetricDeliveryOrders:    float64(m.DeliveryOrders),
		models.MetricActiveListings:    float64(m.ActiveListings),
		models.MetricSellerRating:      m.SellerRating,
		models.MetricAvgReplySeconds:   m.AvgReplySeconds,
		models.MetricFirstListingAt:    float64(m.FirstListingAt),
		models.MetricFirstDealAt:       float64(m.FirstDealAt),
	}
}

func (p predicate) eval(snapshot map[models.MetricKey]float64) bool {
	value, ok := snapshot[p.metric]
	switch p.op {
	case opExists:
		return ok && value != 0
	case opGt:
		return value > p.value
	case opGte:
		return value >= p.value
	case opEq:
		return value == p.value
	default:
		return false
	}
}

func (c condition) eval(snapshot map[models.MetricKey]float64) bool {
	if len(c.predicates) == 0 {
		return true
	}

	if c.match == matchAny {
		for _, p := range c.predicates {
			if p.eval(snapshot) {
				return true
			}
		}
		return false
	}

	for _, p := range c.predicates {
		if !p.eval(snapshot) {
			return false
		}
	}
	return true
}

func (v visibility) allows(mode models.RecapMode) bool {
	switch v {
	case visibilityBoth:
		return true
	case visibilityPublic:
		return mode == models.RecapModePublic
	case visibilityPrivate:
		return mode == models.RecapModePrivate
	default:
		return false
	}
}
