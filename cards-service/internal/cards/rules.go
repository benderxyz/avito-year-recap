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

func (r RuleSet) MetricDefinitions() []models.MetricDefinition {
	return r.metrics
}

type metricSnapshot map[models.MetricKey]float64

func metricsSnapshot(m clients.Metrics, defs []models.MetricDefinition) metricSnapshot {
	snapshot := make(metricSnapshot, len(defs))
	for _, def := range defs {
		value, ok := m.Value(def)
		if !ok {
			continue
		}
		snapshot[models.MetricKey(def.Key)] = value
	}
	return snapshot
}

func (p predicate) eval(snapshot metricSnapshot) bool {
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

func (c condition) eval(snapshot metricSnapshot) bool {
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
