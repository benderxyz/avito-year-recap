package cards

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"cards-service/internal/models"
)

type RuleStore struct {
	db *sql.DB
}

func NewRuleStore(db *sql.DB) *RuleStore {
	return &RuleStore{db: db}
}

func (s *RuleStore) Load(ctx context.Context) (RuleSet, error) {
	badges, err := s.loadBadgeRules(ctx)
	if err != nil {
		return RuleSet{}, err
	}

	stories, err := s.loadStoryRules(ctx)
	if err != nil {
		return RuleSet{}, err
	}

	recommendations, err := s.loadRecommendationRules(ctx)
	if err != nil {
		return RuleSet{}, err
	}

	metrics, err := s.loadMetricDefinitions(ctx)
	if err != nil {
		return RuleSet{}, err
	}

	return RuleSet{
		badges:          badges,
		stories:         stories,
		recommendations: recommendations,
		metrics:         metrics,
	}, nil
}

func (s *RuleStore) loadBadgeRules(ctx context.Context) ([]badgeRule, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, description, icon_url, visibility, metric_key, op, threshold
		FROM badge_rules
		WHERE enabled = true
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, fmt.Errorf("query badge_rules: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var rules []badgeRule
	for rows.Next() {
		var id, title, description, vis, metricKey, op string
		var iconURL sql.NullString
		var threshold float64
		if err := rows.Scan(&id, &title, &description, &iconURL, &vis, &metricKey, &op, &threshold); err != nil {
			return nil, fmt.Errorf("scan badge_rule: %w", err)
		}
		rules = append(rules, makeBadgeRule(id, title, description, iconURL.String, vis, metricKey, op, threshold))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate badge_rules: %w", err)
	}

	return rules, nil
}

func (s *RuleStore) loadStoryRules(ctx context.Context) ([]storyRule, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT visibility, metric_key, op, threshold, payload
		FROM story_rules
		WHERE enabled = true
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, fmt.Errorf("query story_rules: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var rules []storyRule
	for rows.Next() {
		var vis string
		var metricKey, op sql.NullString
		var threshold sql.NullFloat64
		var payload []byte
		if err := rows.Scan(&vis, &metricKey, &op, &threshold, &payload); err != nil {
			return nil, fmt.Errorf("scan story_rule: %w", err)
		}

		rule, err := makeStoryRule(vis, metricKey.String, op.String, threshold.Float64, metricKey.Valid && op.Valid, payload)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate story_rules: %w", err)
	}

	return rules, nil
}

func (s *RuleStore) loadRecommendationRules(ctx context.Context) ([]recommendationRule, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, text, callout, link_label, path, priority, condition
		FROM recommendation_rules
		WHERE enabled = true
		ORDER BY priority DESC, id
	`)
	if err != nil {
		return nil, fmt.Errorf("query recommendation_rules: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var rules []recommendationRule
	for rows.Next() {
		var id, title, text, callout, linkLabel, path string
		var priority int
		var conditionPayload []byte
		if err := rows.Scan(&id, &title, &text, &callout, &linkLabel, &path, &priority, &conditionPayload); err != nil {
			return nil, fmt.Errorf("scan recommendation_rule: %w", err)
		}

		rule, err := makeRecommendationRule(id, title, text, callout, linkLabel, path, priority, conditionPayload)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recommendation_rules: %w", err)
	}

	return rules, nil
}

func (s *RuleStore) loadMetricDefinitions(ctx context.Context) ([]models.MetricDefinition, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT key, value_type, currency, is_public, percentile_key
		FROM metric_definitions
		WHERE enabled = true
		ORDER BY sort_order, key
	`)
	if err != nil {
		return nil, fmt.Errorf("query metric_definitions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var defs []models.MetricDefinition
	for rows.Next() {
		var key, valueType string
		var currency, percentileKey sql.NullString
		var isPublic bool
		if err := rows.Scan(&key, &valueType, &currency, &isPublic, &percentileKey); err != nil {
			return nil, fmt.Errorf("scan metric_definition: %w", err)
		}
		defs = append(defs, models.MetricDefinition{
			Key:           key,
			ValueType:     models.MetricType(valueType),
			Currency:      models.Currency(currency.String),
			IsPublic:      isPublic,
			PercentileKey: percentileKey.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate metric_definitions: %w", err)
	}

	return defs, nil
}

func makeBadgeRule(id, title, description, iconURL, vis, metricKey, op string, threshold float64) badgeRule {
	return badgeRule{
		badge: models.Badge{
			ID:          id,
			Title:       title,
			Description: description,
			IconURL:     iconURL,
		},
		visibility: visibility(vis),
		when: condition{predicates: []predicate{
			{metric: models.MetricKey(metricKey), op: predicateOp(op), value: threshold},
		}},
	}
}

func makeStoryRule(vis, metricKey, op string, threshold float64, hasPredicate bool, payload []byte) (storyRule, error) {
	var scene map[string]any
	if err := json.Unmarshal(payload, &scene); err != nil {
		return storyRule{}, fmt.Errorf("unmarshal story payload: %w", err)
	}

	rule := storyRule{
		visibility: visibility(vis),
		scene:      scene,
	}
	if hasPredicate {
		rule.when = condition{predicates: []predicate{
			{metric: models.MetricKey(metricKey), op: predicateOp(op), value: threshold},
		}}
	}

	return rule, nil
}

func makeRecommendationRule(id, title, text, callout, linkLabel, path string, priority int, conditionPayload []byte) (recommendationRule, error) {
	when, err := parseCondition(conditionPayload)
	if err != nil {
		return recommendationRule{}, err
	}

	return recommendationRule{
		id:        id,
		title:     title,
		text:      text,
		callout:   callout,
		linkLabel: linkLabel,
		path:      path,
		priority:  priority,
		when:      when,
	}, nil
}

func parseCondition(payload []byte) (condition, error) {
	if len(payload) == 0 {
		return condition{}, nil
	}

	var dto struct {
		Match      string `json:"match"`
		Predicates []struct {
			Metric string  `json:"metric"`
			Op     string  `json:"op"`
			Value  float64 `json:"value"`
		} `json:"predicates"`
	}
	if err := json.Unmarshal(payload, &dto); err != nil {
		return condition{}, fmt.Errorf("unmarshal condition: %w", err)
	}

	result := condition{match: matchMode(dto.Match)}
	for _, p := range dto.Predicates {
		result.predicates = append(result.predicates, predicate{
			metric: models.MetricKey(p.Metric),
			op:     predicateOp(p.Op),
			value:  p.Value,
		})
	}

	return result, nil
}
