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

	return RuleSet{badges: badges, stories: stories}, nil
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
