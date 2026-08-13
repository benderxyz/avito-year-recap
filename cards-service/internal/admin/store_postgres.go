package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

type rowScanner interface {
	Scan(dest ...any) error
}

type conditions struct {
	clauses []string
	args    []any
}

func (c *conditions) add(format string, value any) {
	c.args = append(c.args, value)
	c.clauses = append(c.clauses, fmt.Sprintf(format, "$"+strconv.Itoa(len(c.args))))
}

func (c *conditions) where() string {
	if len(c.clauses) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(c.clauses, " AND ")
}

func listQuery(columns, table, where, order string) string {
	return "SELECT " + columns + " FROM " + table + where + " ORDER BY " + order
}

func timestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func nullableString(value *string) any {
	if value == nil || *value == "" {
		return nil
	}
	return *value
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func nullableFloat(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func stringPointer(value sql.NullString) *string {
	if !value.Valid || value.String == "" {
		return nil
	}
	result := value.String
	return &result
}

func floatPointer(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	result := value.Float64
	return &result
}

func isUniqueViolation(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return string(pgErr.Code) == uniqueViolationCode
	}
	return false
}

const metricColumns = `key, value_type, currency, is_public, percentile_key, comparison_min_percentile,
	COALESCE(source_key, key), source_field, include_in_llm, sort_order, enabled, created_at, updated_at`

func scanMetric(row rowScanner) (MetricDefinition, error) {
	var (
		metric                          MetricDefinition
		currency, percentileKey         sql.NullString
		comparisonMin                   sql.NullFloat64
		valueType, sourceField          string
		sourceKey                       string
		createdAt, updatedAt            time.Time
		isPublic, includeInLlm, enabled bool
	)

	if err := row.Scan(&metric.Key, &valueType, &currency, &isPublic, &percentileKey, &comparisonMin,
		&sourceKey, &sourceField, &includeInLlm, &metric.SortOrder, &enabled, &createdAt, &updatedAt); err != nil {
		return MetricDefinition{}, err
	}

	metric.ValueType = MetricValueType(valueType)
	metric.Currency = stringPointer(currency)
	metric.IsPublic = isPublic
	metric.PercentileKey = stringPointer(percentileKey)
	metric.ComparisonMinPercentile = floatPointer(comparisonMin)
	metric.SourceKey = sourceKey
	metric.SourceField = MetricSourceField(sourceField)
	metric.IncludeInLlm = includeInLlm
	metric.Enabled = enabled
	metric.CreatedAt = timestamp(createdAt)
	metric.UpdatedAt = timestamp(updatedAt)

	return metric, nil
}

func (s *PostgresStore) ListMetrics(ctx context.Context, filter MetricFilter) ([]MetricDefinition, error) {
	var where conditions
	if filter.Enabled != nil {
		where.add("enabled = %s", *filter.Enabled)
	}
	if filter.IsPublic != nil {
		where.add("is_public = %s", *filter.IsPublic)
	}
	if filter.IncludeInLlm != nil {
		where.add("include_in_llm = %s", *filter.IncludeInLlm)
	}
	if filter.ValueType != "" {
		where.add("value_type = %s", string(filter.ValueType))
	}
	if filter.SourceField != "" {
		where.add("source_field = %s", string(filter.SourceField))
	}
	if filter.Search != "" {
		where.add("key ILIKE %s", "%"+filter.Search+"%")
	}

	query := listQuery(metricColumns, "metric_definitions", where.where(), "sort_order, key")

	rows, err := s.db.QueryContext(ctx, query, where.args...)
	if err != nil {
		return nil, fmt.Errorf("query metric_definitions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	metrics := make([]MetricDefinition, 0)
	for rows.Next() {
		metric, err := scanMetric(rows)
		if err != nil {
			return nil, fmt.Errorf("scan metric_definition: %w", err)
		}
		metrics = append(metrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate metric_definitions: %w", err)
	}

	return metrics, nil
}

func (s *PostgresStore) GetMetric(ctx context.Context, key string) (MetricDefinition, error) {
	row := s.db.QueryRowContext(ctx, "SELECT "+metricColumns+" FROM metric_definitions WHERE key = $1", key)

	metric, err := scanMetric(row)
	if errors.Is(err, sql.ErrNoRows) {
		return MetricDefinition{}, ErrNotFound
	}
	if err != nil {
		return MetricDefinition{}, fmt.Errorf("get metric_definition: %w", err)
	}

	return metric, nil
}

func metricArgs(key string, write MetricWrite) []any {
	sourceKey := write.SourceKey
	if sourceKey == "" {
		sourceKey = key
	}
	sourceField := write.SourceField
	if sourceField == "" {
		sourceField = SourceFieldValue
	}

	return []any{
		key,
		string(write.ValueType),
		nullableString(write.Currency),
		write.IsPublic,
		nullableString(write.PercentileKey),
		nullableFloat(write.ComparisonMinPercentile),
		sourceKey,
		string(sourceField),
		write.IncludeInLlm,
		write.SortOrder,
		write.Enabled,
	}
}

func (s *PostgresStore) CreateMetric(ctx context.Context, key string, write MetricWrite) (MetricDefinition, error) {
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO metric_definitions (key, value_type, currency, is_public, percentile_key,
			comparison_min_percentile, source_key, source_field, include_in_llm, sort_order, enabled,
			created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now())
		RETURNING `+metricColumns, metricArgs(key, write)...)

	metric, err := scanMetric(row)
	if isUniqueViolation(err) {
		return MetricDefinition{}, ErrAlreadyExists
	}
	if err != nil {
		return MetricDefinition{}, fmt.Errorf("insert metric_definition: %w", err)
	}

	return metric, nil
}

func (s *PostgresStore) UpdateMetric(ctx context.Context, key string, write MetricWrite) (MetricDefinition, error) {
	row := s.db.QueryRowContext(ctx, `
		UPDATE metric_definitions
		SET value_type = $2, currency = $3, is_public = $4, percentile_key = $5,
			comparison_min_percentile = $6, source_key = $7, source_field = $8,
			include_in_llm = $9, sort_order = $10, enabled = $11, updated_at = now()
		WHERE key = $1
		RETURNING `+metricColumns, metricArgs(key, write)...)

	metric, err := scanMetric(row)
	if errors.Is(err, sql.ErrNoRows) {
		return MetricDefinition{}, ErrNotFound
	}
	if err != nil {
		return MetricDefinition{}, fmt.Errorf("update metric_definition: %w", err)
	}

	return metric, nil
}

func (s *PostgresStore) DeleteMetric(ctx context.Context, key string) error {
	return s.deleteByID(ctx, "DELETE FROM metric_definitions WHERE key = $1", key)
}

func (s *PostgresStore) MetricReferences(ctx context.Context, key string) ([]string, error) {
	return s.references(ctx, `
		SELECT 'badge:' || id FROM badge_rules WHERE metric_key = $1
		UNION ALL
		SELECT 'story:' || id FROM story_rules
		WHERE metric_key = $1
		   OR payload->>'value' = $1
		   OR payload->>'percentile' = $1
		   OR payload->'blocks' @> jsonb_build_array(jsonb_build_object('value', $1::text))
		UNION ALL
		SELECT 'recommendation:' || id FROM recommendation_rules
		WHERE condition->'predicates' @> jsonb_build_array(jsonb_build_object('metric', $1::text))
		UNION ALL
		SELECT 'metric:' || key FROM metric_definitions
		WHERE key <> $1 AND (percentile_key = $1 OR source_key = $1)
	`, key)
}

const badgeColumns = `id, title, description, icon_url, visibility, metric_key, op, threshold,
	sort_order, enabled, created_at, updated_at`

func scanBadge(row rowScanner) (BadgeRule, error) {
	var (
		badge                string
		title, description   string
		iconURL              sql.NullString
		vis, metricKey, op   string
		threshold            float64
		sortOrder            int
		enabled              bool
		createdAt, updatedAt time.Time
	)

	if err := row.Scan(&badge, &title, &description, &iconURL, &vis, &metricKey, &op, &threshold,
		&sortOrder, &enabled, &createdAt, &updatedAt); err != nil {
		return BadgeRule{}, err
	}

	return BadgeRule{
		ID:          badge,
		Title:       title,
		Description: description,
		IconURL:     stringPointer(iconURL),
		Visibility:  Visibility(vis),
		When:        SimpleWhen{Metric: metricKey, Op: PredicateOp(op), Value: threshold},
		SortOrder:   sortOrder,
		Enabled:     enabled,
		CreatedAt:   timestamp(createdAt),
		UpdatedAt:   timestamp(updatedAt),
	}, nil
}

func (s *PostgresStore) ListBadges(ctx context.Context, filter BadgeFilter) ([]BadgeRule, error) {
	var where conditions
	if filter.Enabled != nil {
		where.add("enabled = %s", *filter.Enabled)
	}
	if filter.Visibility != "" {
		where.add("visibility = %s", string(filter.Visibility))
	}
	if filter.Metric != "" {
		where.add("metric_key = %s", filter.Metric)
	}
	if filter.Search != "" {
		where.add("(id ILIKE %[1]s OR title ILIKE %[1]s OR description ILIKE %[1]s)", "%"+filter.Search+"%")
	}

	query := listQuery(badgeColumns, "badge_rules", where.where(), "sort_order, id")

	rows, err := s.db.QueryContext(ctx, query, where.args...)
	if err != nil {
		return nil, fmt.Errorf("query badge_rules: %w", err)
	}
	defer func() { _ = rows.Close() }()

	badges := make([]BadgeRule, 0)
	for rows.Next() {
		badge, err := scanBadge(rows)
		if err != nil {
			return nil, fmt.Errorf("scan badge_rule: %w", err)
		}
		badges = append(badges, badge)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate badge_rules: %w", err)
	}

	return badges, nil
}

func (s *PostgresStore) GetBadge(ctx context.Context, id string) (BadgeRule, error) {
	row := s.db.QueryRowContext(ctx, "SELECT "+badgeColumns+" FROM badge_rules WHERE id = $1", id)

	badge, err := scanBadge(row)
	if errors.Is(err, sql.ErrNoRows) {
		return BadgeRule{}, ErrNotFound
	}
	if err != nil {
		return BadgeRule{}, fmt.Errorf("get badge_rule: %w", err)
	}

	return badge, nil
}

func badgeArgs(id string, write BadgeWrite) []any {
	return []any{
		id,
		write.Title,
		write.Description,
		nullableString(write.IconURL),
		string(write.Visibility),
		write.When.Metric,
		string(write.When.Op),
		write.When.Value,
		write.SortOrder,
		write.Enabled,
	}
}

func (s *PostgresStore) CreateBadge(ctx context.Context, id string, write BadgeWrite) (BadgeRule, error) {
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO badge_rules (id, title, description, icon_url, visibility, metric_key, op,
			threshold, sort_order, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now())
		RETURNING `+badgeColumns, badgeArgs(id, write)...)

	badge, err := scanBadge(row)
	if isUniqueViolation(err) {
		return BadgeRule{}, ErrAlreadyExists
	}
	if err != nil {
		return BadgeRule{}, fmt.Errorf("insert badge_rule: %w", err)
	}

	return badge, nil
}

func (s *PostgresStore) UpdateBadge(ctx context.Context, id string, write BadgeWrite) (BadgeRule, error) {
	row := s.db.QueryRowContext(ctx, `
		UPDATE badge_rules
		SET title = $2, description = $3, icon_url = $4, visibility = $5, metric_key = $6,
			op = $7, threshold = $8, sort_order = $9, enabled = $10, updated_at = now()
		WHERE id = $1
		RETURNING `+badgeColumns, badgeArgs(id, write)...)

	badge, err := scanBadge(row)
	if errors.Is(err, sql.ErrNoRows) {
		return BadgeRule{}, ErrNotFound
	}
	if err != nil {
		return BadgeRule{}, fmt.Errorf("update badge_rule: %w", err)
	}

	return badge, nil
}

func (s *PostgresStore) DeleteBadge(ctx context.Context, id string) error {
	return s.deleteByID(ctx, "DELETE FROM badge_rules WHERE id = $1", id)
}

func (s *PostgresStore) BadgeReferences(ctx context.Context, id string) ([]string, error) {
	return s.references(ctx, `
		SELECT 'story:' || id FROM story_rules WHERE payload->>'badgeId' = $1
	`, id)
}

const storyColumns = `id, scene_type, visibility, metric_key, op, threshold, payload,
	sort_order, enabled, created_at, updated_at`

func scanStory(row rowScanner) (StoryRule, error) {
	var (
		id, sceneType, vis   string
		metricKey, op        sql.NullString
		threshold            sql.NullFloat64
		payload              []byte
		sortOrder            int
		enabled              bool
		createdAt, updatedAt time.Time
	)

	if err := row.Scan(&id, &sceneType, &vis, &metricKey, &op, &threshold, &payload,
		&sortOrder, &enabled, &createdAt, &updatedAt); err != nil {
		return StoryRule{}, err
	}

	scene := map[string]any{}
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &scene); err != nil {
			return StoryRule{}, fmt.Errorf("unmarshal story payload: %w", err)
		}
	}

	story := StoryRule{
		ID:         id,
		SceneType:  SceneType(sceneType),
		Visibility: Visibility(vis),
		Payload:    scene,
		SortOrder:  sortOrder,
		Enabled:    enabled,
		CreatedAt:  timestamp(createdAt),
		UpdatedAt:  timestamp(updatedAt),
	}

	if metricKey.Valid && op.Valid {
		story.When = &SimpleWhen{
			Metric: metricKey.String,
			Op:     PredicateOp(op.String),
			Value:  threshold.Float64,
		}
	}

	return story, nil
}

func (s *PostgresStore) ListStories(ctx context.Context, filter StoryFilter) ([]StoryRule, error) {
	var where conditions
	if filter.Enabled != nil {
		where.add("enabled = %s", *filter.Enabled)
	}
	if filter.Visibility != "" {
		where.add("visibility = %s", string(filter.Visibility))
	}
	if filter.SceneType != "" {
		where.add("scene_type = %s", string(filter.SceneType))
	}
	if filter.Metric != "" {
		where.add("metric_key = %s", filter.Metric)
	}
	if filter.Search != "" {
		where.add("(id ILIKE %[1]s OR payload::text ILIKE %[1]s)", "%"+filter.Search+"%")
	}

	query := listQuery(storyColumns, "story_rules", where.where(), "sort_order, id")

	rows, err := s.db.QueryContext(ctx, query, where.args...)
	if err != nil {
		return nil, fmt.Errorf("query story_rules: %w", err)
	}
	defer func() { _ = rows.Close() }()

	stories := make([]StoryRule, 0)
	for rows.Next() {
		story, err := scanStory(rows)
		if err != nil {
			return nil, fmt.Errorf("scan story_rule: %w", err)
		}
		stories = append(stories, story)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate story_rules: %w", err)
	}

	return stories, nil
}

func (s *PostgresStore) GetStory(ctx context.Context, id string) (StoryRule, error) {
	row := s.db.QueryRowContext(ctx, "SELECT "+storyColumns+" FROM story_rules WHERE id = $1", id)

	story, err := scanStory(row)
	if errors.Is(err, sql.ErrNoRows) {
		return StoryRule{}, ErrNotFound
	}
	if err != nil {
		return StoryRule{}, fmt.Errorf("get story_rule: %w", err)
	}

	return story, nil
}

func storyArgs(id string, write StoryWrite) ([]any, error) {
	payload, err := json.Marshal(write.Payload)
	if err != nil {
		return nil, fmt.Errorf("marshal story payload: %w", err)
	}

	var metricKey, op any
	var threshold any
	if write.When != nil {
		metricKey = write.When.Metric
		op = string(write.When.Op)
		threshold = write.When.Value
	}

	return []any{
		id,
		string(write.SceneType),
		string(write.Visibility),
		metricKey,
		op,
		threshold,
		payload,
		write.SortOrder,
		write.Enabled,
	}, nil
}

func (s *PostgresStore) CreateStory(ctx context.Context, id string, write StoryWrite) (StoryRule, error) {
	args, err := storyArgs(id, write)
	if err != nil {
		return StoryRule{}, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO story_rules (id, scene_type, visibility, metric_key, op, threshold, payload,
			sort_order, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now())
		RETURNING `+storyColumns, args...)

	story, err := scanStory(row)
	if isUniqueViolation(err) {
		return StoryRule{}, ErrAlreadyExists
	}
	if err != nil {
		return StoryRule{}, fmt.Errorf("insert story_rule: %w", err)
	}

	return story, nil
}

func (s *PostgresStore) UpdateStory(ctx context.Context, id string, write StoryWrite) (StoryRule, error) {
	args, err := storyArgs(id, write)
	if err != nil {
		return StoryRule{}, err
	}

	row := s.db.QueryRowContext(ctx, `
		UPDATE story_rules
		SET scene_type = $2, visibility = $3, metric_key = $4, op = $5, threshold = $6,
			payload = $7, sort_order = $8, enabled = $9, updated_at = now()
		WHERE id = $1
		RETURNING `+storyColumns, args...)

	story, err := scanStory(row)
	if errors.Is(err, sql.ErrNoRows) {
		return StoryRule{}, ErrNotFound
	}
	if err != nil {
		return StoryRule{}, fmt.Errorf("update story_rule: %w", err)
	}

	return story, nil
}

func (s *PostgresStore) DeleteStory(ctx context.Context, id string) error {
	return s.deleteByID(ctx, "DELETE FROM story_rules WHERE id = $1", id)
}

const recommendationColumns = `id, title, text, callout, link_label, path, priority, condition,
	enabled, created_at, updated_at`

func scanRecommendation(row rowScanner) (RecommendationRule, error) {
	var (
		id, title, text          string
		callout, linkLabel, path string
		priority                 int
		conditionPayload         []byte
		enabled                  bool
		createdAt, updatedAt     time.Time
	)

	when := GroupWhen{Predicates: []Predicate{}}

	if err := row.Scan(&id, &title, &text, &callout, &linkLabel, &path, &priority,
		&conditionPayload, &enabled, &createdAt, &updatedAt); err != nil {
		return RecommendationRule{}, err
	}

	if len(conditionPayload) > 0 {
		if err := json.Unmarshal(conditionPayload, &when); err != nil {
			return RecommendationRule{}, fmt.Errorf("unmarshal recommendation condition: %w", err)
		}
	}
	if when.Predicates == nil {
		when.Predicates = []Predicate{}
	}

	return RecommendationRule{
		ID:        id,
		Title:     title,
		Text:      text,
		Callout:   callout,
		LinkLabel: linkLabel,
		Path:      path,
		Priority:  priority,
		When:      when,
		Enabled:   enabled,
		CreatedAt: timestamp(createdAt),
		UpdatedAt: timestamp(updatedAt),
	}, nil
}

func (s *PostgresStore) ListRecommendations(ctx context.Context, filter RecommendationFilter) ([]RecommendationRule, error) {
	var where conditions
	if filter.Enabled != nil {
		where.add("enabled = %s", *filter.Enabled)
	}
	if filter.MinPriority != nil {
		where.add("priority >= %s", *filter.MinPriority)
	}
	if filter.Metric != "" {
		where.add("condition->'predicates' @> jsonb_build_array(jsonb_build_object('metric', %s::text))", filter.Metric)
	}
	if filter.Search != "" {
		where.add("(id ILIKE %[1]s OR title ILIKE %[1]s OR text ILIKE %[1]s)", "%"+filter.Search+"%")
	}

	query := listQuery(recommendationColumns, "recommendation_rules", where.where(), "priority DESC, id")

	rows, err := s.db.QueryContext(ctx, query, where.args...)
	if err != nil {
		return nil, fmt.Errorf("query recommendation_rules: %w", err)
	}
	defer func() { _ = rows.Close() }()

	recommendations := make([]RecommendationRule, 0)
	for rows.Next() {
		recommendation, err := scanRecommendation(rows)
		if err != nil {
			return nil, fmt.Errorf("scan recommendation_rule: %w", err)
		}
		recommendations = append(recommendations, recommendation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recommendation_rules: %w", err)
	}

	return recommendations, nil
}

func (s *PostgresStore) GetRecommendation(ctx context.Context, id string) (RecommendationRule, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT "+recommendationColumns+" FROM recommendation_rules WHERE id = $1", id)

	recommendation, err := scanRecommendation(row)
	if errors.Is(err, sql.ErrNoRows) {
		return RecommendationRule{}, ErrNotFound
	}
	if err != nil {
		return RecommendationRule{}, fmt.Errorf("get recommendation_rule: %w", err)
	}

	return recommendation, nil
}

func recommendationArgs(id string, write RecommendationWrite) ([]any, error) {
	when := write.When
	if when.Predicates == nil {
		when.Predicates = []Predicate{}
	}

	condition, err := json.Marshal(when)
	if err != nil {
		return nil, fmt.Errorf("marshal recommendation condition: %w", err)
	}

	return []any{
		id,
		write.Title,
		write.Text,
		write.Callout,
		write.LinkLabel,
		write.Path,
		write.Priority,
		condition,
		write.Enabled,
	}, nil
}

func (s *PostgresStore) CreateRecommendation(ctx context.Context, id string, write RecommendationWrite) (RecommendationRule, error) {
	args, err := recommendationArgs(id, write)
	if err != nil {
		return RecommendationRule{}, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO recommendation_rules (id, title, text, callout, link_label, path, priority,
			condition, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now())
		RETURNING `+recommendationColumns, args...)

	recommendation, err := scanRecommendation(row)
	if isUniqueViolation(err) {
		return RecommendationRule{}, ErrAlreadyExists
	}
	if err != nil {
		return RecommendationRule{}, fmt.Errorf("insert recommendation_rule: %w", err)
	}

	return recommendation, nil
}

func (s *PostgresStore) UpdateRecommendation(ctx context.Context, id string, write RecommendationWrite) (RecommendationRule, error) {
	args, err := recommendationArgs(id, write)
	if err != nil {
		return RecommendationRule{}, err
	}

	row := s.db.QueryRowContext(ctx, `
		UPDATE recommendation_rules
		SET title = $2, text = $3, callout = $4, link_label = $5, path = $6, priority = $7,
			condition = $8, enabled = $9, updated_at = now()
		WHERE id = $1
		RETURNING `+recommendationColumns, args...)

	recommendation, err := scanRecommendation(row)
	if errors.Is(err, sql.ErrNoRows) {
		return RecommendationRule{}, ErrNotFound
	}
	if err != nil {
		return RecommendationRule{}, fmt.Errorf("update recommendation_rule: %w", err)
	}

	return recommendation, nil
}

func (s *PostgresStore) DeleteRecommendation(ctx context.Context, id string) error {
	return s.deleteByID(ctx, "DELETE FROM recommendation_rules WHERE id = $1", id)
}

func (s *PostgresStore) deleteByID(ctx context.Context, query, id string) error {
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresStore) references(ctx context.Context, query, id string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query references: %w", err)
	}
	defer func() { _ = rows.Close() }()

	refs := make([]string, 0)
	for rows.Next() {
		var ref string
		if err := rows.Scan(&ref); err != nil {
			return nil, fmt.Errorf("scan reference: %w", err)
		}
		refs = append(refs, ref)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate references: %w", err)
	}

	return refs, nil
}
