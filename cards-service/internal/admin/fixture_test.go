package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"cards-service/internal/cards"
)

const testToken = "test-admin-token"

type memoryStore struct {
	metrics         map[string]MetricDefinition
	badges          map[string]BadgeRule
	stories         map[string]StoryRule
	recommendations map[string]RecommendationRule
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		metrics:         map[string]MetricDefinition{},
		badges:          map[string]BadgeRule{},
		stories:         map[string]StoryRule{},
		recommendations: map[string]RecommendationRule{},
	}
}

const (
	testCreatedAt = "2026-01-01T00:00:00Z"
	testUpdatedAt = "2026-02-02T00:00:00Z"
)

func (s *memoryStore) ListMetrics(_ context.Context, filter MetricFilter) ([]MetricDefinition, error) {
	items := make([]MetricDefinition, 0, len(s.metrics))
	for _, metric := range s.metrics {
		switch {
		case filter.Enabled != nil && metric.Enabled != *filter.Enabled:
		case filter.IsPublic != nil && metric.IsPublic != *filter.IsPublic:
		case filter.IncludeInLlm != nil && metric.IncludeInLlm != *filter.IncludeInLlm:
		case filter.ValueType != "" && metric.ValueType != filter.ValueType:
		case filter.SourceField != "" && metric.SourceField != filter.SourceField:
		case filter.Search != "" && !strings.Contains(metric.Key, filter.Search):
		default:
			items = append(items, metric)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].SortOrder != items[j].SortOrder {
			return items[i].SortOrder < items[j].SortOrder
		}
		return items[i].Key < items[j].Key
	})

	return items, nil
}

func (s *memoryStore) GetMetric(_ context.Context, key string) (MetricDefinition, error) {
	metric, ok := s.metrics[key]
	if !ok {
		return MetricDefinition{}, ErrNotFound
	}
	return metric, nil
}

func metricFromWrite(key string, write MetricWrite, createdAt string) MetricDefinition {
	sourceKey := write.SourceKey
	if sourceKey == "" {
		sourceKey = key
	}
	sourceField := write.SourceField
	if sourceField == "" {
		sourceField = SourceFieldValue
	}

	return MetricDefinition{
		Key:                     key,
		ValueType:               write.ValueType,
		Currency:                write.Currency,
		IsPublic:                write.IsPublic,
		PercentileKey:           write.PercentileKey,
		ComparisonMinPercentile: write.ComparisonMinPercentile,
		SourceKey:               sourceKey,
		SourceField:             sourceField,
		IncludeInLlm:            write.IncludeInLlm,
		SortOrder:               write.SortOrder,
		Enabled:                 write.Enabled,
		CreatedAt:               createdAt,
		UpdatedAt:               testUpdatedAt,
	}
}

func (s *memoryStore) CreateMetric(_ context.Context, key string, write MetricWrite) (MetricDefinition, error) {
	if _, exists := s.metrics[key]; exists {
		return MetricDefinition{}, ErrAlreadyExists
	}

	metric := metricFromWrite(key, write, testCreatedAt)
	s.metrics[key] = metric

	return metric, nil
}

func (s *memoryStore) UpdateMetric(_ context.Context, key string, write MetricWrite) (MetricDefinition, error) {
	existing, ok := s.metrics[key]
	if !ok {
		return MetricDefinition{}, ErrNotFound
	}

	metric := metricFromWrite(key, write, existing.CreatedAt)
	s.metrics[key] = metric

	return metric, nil
}

func (s *memoryStore) DeleteMetric(_ context.Context, key string) error {
	if _, ok := s.metrics[key]; !ok {
		return ErrNotFound
	}
	delete(s.metrics, key)
	return nil
}

func (s *memoryStore) MetricReferences(_ context.Context, key string) ([]string, error) {
	refs := make([]string, 0)

	for _, badge := range s.badges {
		if badge.When.Metric == key {
			refs = append(refs, "badge:"+badge.ID)
		}
	}
	for _, story := range s.stories {
		value, _ := story.Payload["value"].(string)
		if (story.When != nil && story.When.Metric == key) || value == key {
			refs = append(refs, "story:"+story.ID)
		}
	}
	for _, recommendation := range s.recommendations {
		for _, predicate := range recommendation.When.Predicates {
			if predicate.Metric == key {
				refs = append(refs, "recommendation:"+recommendation.ID)
				break
			}
		}
	}
	for _, metric := range s.metrics {
		if metric.Key != key && (derefString(metric.PercentileKey) == key || metric.SourceKey == key) {
			refs = append(refs, "metric:"+metric.Key)
		}
	}

	sort.Strings(refs)

	return refs, nil
}

func (s *memoryStore) ListBadges(_ context.Context, filter BadgeFilter) ([]BadgeRule, error) {
	items := make([]BadgeRule, 0, len(s.badges))
	for _, badge := range s.badges {
		switch {
		case filter.Enabled != nil && badge.Enabled != *filter.Enabled:
		case filter.Visibility != "" && badge.Visibility != filter.Visibility:
		case filter.Metric != "" && badge.When.Metric != filter.Metric:
		case filter.Search != "" && !strings.Contains(badge.ID+badge.Title, filter.Search):
		default:
			items = append(items, badge)
		}
	}

	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })

	return items, nil
}

func (s *memoryStore) GetBadge(_ context.Context, id string) (BadgeRule, error) {
	badge, ok := s.badges[id]
	if !ok {
		return BadgeRule{}, ErrNotFound
	}
	return badge, nil
}

func badgeFromWrite(id string, write BadgeWrite, createdAt string) BadgeRule {
	return BadgeRule{
		ID:          id,
		Title:       write.Title,
		Description: write.Description,
		IconURL:     write.IconURL,
		Visibility:  write.Visibility,
		When:        write.When,
		SortOrder:   write.SortOrder,
		Enabled:     write.Enabled,
		CreatedAt:   createdAt,
		UpdatedAt:   testUpdatedAt,
	}
}

func (s *memoryStore) CreateBadge(_ context.Context, id string, write BadgeWrite) (BadgeRule, error) {
	if _, exists := s.badges[id]; exists {
		return BadgeRule{}, ErrAlreadyExists
	}

	badge := badgeFromWrite(id, write, testCreatedAt)
	s.badges[id] = badge

	return badge, nil
}

func (s *memoryStore) UpdateBadge(_ context.Context, id string, write BadgeWrite) (BadgeRule, error) {
	existing, ok := s.badges[id]
	if !ok {
		return BadgeRule{}, ErrNotFound
	}

	badge := badgeFromWrite(id, write, existing.CreatedAt)
	s.badges[id] = badge

	return badge, nil
}

func (s *memoryStore) DeleteBadge(_ context.Context, id string) error {
	if _, ok := s.badges[id]; !ok {
		return ErrNotFound
	}
	delete(s.badges, id)
	return nil
}

func (s *memoryStore) BadgeReferences(_ context.Context, id string) ([]string, error) {
	refs := make([]string, 0)
	for _, story := range s.stories {
		if badgeID, _ := story.Payload["badgeId"].(string); badgeID == id {
			refs = append(refs, "story:"+story.ID)
		}
	}
	sort.Strings(refs)
	return refs, nil
}

func (s *memoryStore) ListStories(_ context.Context, filter StoryFilter) ([]StoryRule, error) {
	items := make([]StoryRule, 0, len(s.stories))
	for _, story := range s.stories {
		switch {
		case filter.Enabled != nil && story.Enabled != *filter.Enabled:
		case filter.Visibility != "" && story.Visibility != filter.Visibility:
		case filter.SceneType != "" && story.SceneType != filter.SceneType:
		case filter.Search != "" && !strings.Contains(story.ID, filter.Search):
		default:
			items = append(items, story)
		}
	}

	sort.Slice(items, func(i, j int) bool { return items[i].SortOrder < items[j].SortOrder })

	return items, nil
}

func (s *memoryStore) GetStory(_ context.Context, id string) (StoryRule, error) {
	story, ok := s.stories[id]
	if !ok {
		return StoryRule{}, ErrNotFound
	}
	return story, nil
}

func storyFromWrite(id string, write StoryWrite, createdAt string) StoryRule {
	return StoryRule{
		ID:         id,
		SceneType:  write.SceneType,
		Visibility: write.Visibility,
		When:       write.When,
		Payload:    write.Payload,
		SortOrder:  write.SortOrder,
		Enabled:    write.Enabled,
		CreatedAt:  createdAt,
		UpdatedAt:  testUpdatedAt,
	}
}

func (s *memoryStore) CreateStory(_ context.Context, id string, write StoryWrite) (StoryRule, error) {
	if _, exists := s.stories[id]; exists {
		return StoryRule{}, ErrAlreadyExists
	}

	story := storyFromWrite(id, write, testCreatedAt)
	s.stories[id] = story

	return story, nil
}

func (s *memoryStore) UpdateStory(_ context.Context, id string, write StoryWrite) (StoryRule, error) {
	existing, ok := s.stories[id]
	if !ok {
		return StoryRule{}, ErrNotFound
	}

	story := storyFromWrite(id, write, existing.CreatedAt)
	s.stories[id] = story

	return story, nil
}

func (s *memoryStore) DeleteStory(_ context.Context, id string) error {
	if _, ok := s.stories[id]; !ok {
		return ErrNotFound
	}
	delete(s.stories, id)
	return nil
}

func (s *memoryStore) ListRecommendations(_ context.Context, filter RecommendationFilter) ([]RecommendationRule, error) {
	items := make([]RecommendationRule, 0, len(s.recommendations))
	for _, recommendation := range s.recommendations {
		switch {
		case filter.Enabled != nil && recommendation.Enabled != *filter.Enabled:
		case filter.MinPriority != nil && recommendation.Priority < *filter.MinPriority:
		case filter.Search != "" && !strings.Contains(recommendation.ID+recommendation.Title, filter.Search):
		default:
			items = append(items, recommendation)
		}
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Priority > items[j].Priority })

	return items, nil
}

func (s *memoryStore) GetRecommendation(_ context.Context, id string) (RecommendationRule, error) {
	recommendation, ok := s.recommendations[id]
	if !ok {
		return RecommendationRule{}, ErrNotFound
	}
	return recommendation, nil
}

func recommendationFromWrite(id string, write RecommendationWrite, createdAt string) RecommendationRule {
	when := write.When
	if when.Predicates == nil {
		when.Predicates = []Predicate{}
	}

	return RecommendationRule{
		ID:        id,
		Title:     write.Title,
		Text:      write.Text,
		Callout:   write.Callout,
		LinkLabel: write.LinkLabel,
		Path:      write.Path,
		Priority:  write.Priority,
		When:      when,
		Enabled:   write.Enabled,
		CreatedAt: createdAt,
		UpdatedAt: testUpdatedAt,
	}
}

func (s *memoryStore) CreateRecommendation(_ context.Context, id string, write RecommendationWrite) (RecommendationRule, error) {
	if _, exists := s.recommendations[id]; exists {
		return RecommendationRule{}, ErrAlreadyExists
	}

	recommendation := recommendationFromWrite(id, write, testCreatedAt)
	s.recommendations[id] = recommendation

	return recommendation, nil
}

func (s *memoryStore) UpdateRecommendation(_ context.Context, id string, write RecommendationWrite) (RecommendationRule, error) {
	existing, ok := s.recommendations[id]
	if !ok {
		return RecommendationRule{}, ErrNotFound
	}

	recommendation := recommendationFromWrite(id, write, existing.CreatedAt)
	s.recommendations[id] = recommendation

	return recommendation, nil
}

func (s *memoryStore) DeleteRecommendation(_ context.Context, id string) error {
	if _, ok := s.recommendations[id]; !ok {
		return ErrNotFound
	}
	delete(s.recommendations, id)
	return nil
}

type stubProvider struct {
	invalidated int
	ruleSet     cards.RuleSet
	err         error
}

func (p *stubProvider) Get(context.Context) (cards.RuleSet, error) {
	return p.ruleSet, p.err
}

func (p *stubProvider) Invalidate() {
	p.invalidated++
}

func newTestServer(t *testing.T, store Store, provider RuleProvider) http.Handler {
	t.Helper()

	mux := http.NewServeMux()
	RegisterRoutes(mux, NewHandler(Options{
		Token:          testToken,
		Store:          store,
		Rules:          provider,
		ShareBaseURL:   "http://localhost:3000",
		ProductBaseURL: "https://www.avito.ru",
	}))

	return mux
}

func request(t *testing.T, handler http.Handler, method, target, token string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = strings.NewReader(string(encoded))
	} else {
		reader = strings.NewReader("")
	}

	req := httptest.NewRequestWithContext(context.Background(), method, target, reader)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	return recorder
}

func decodeResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()

	var payload T
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}

	return payload
}

func newBadgeBody(id, metric string) BadgeCreate {
	return BadgeCreate{
		ID: id,
		BadgeWrite: BadgeWrite{
			Title:       "Разговорчивый",
			Description: "Вы отправили больше 1000 сообщений",
			Visibility:  VisibilityBoth,
			When:        SimpleWhen{Metric: metric, Op: OpGt, Value: 1000},
			SortOrder:   10,
			Enabled:     true,
		},
	}
}

func newStoryBody(id string) StoryCreate {
	return StoryCreate{
		ID: id,
		StoryWrite: StoryWrite{
			SceneType:  SceneIntro,
			Visibility: VisibilityBoth,
			Payload: map[string]any{
				"id":    id,
				"type":  string(SceneIntro),
				"title": "Ваш год на Авито",
			},
			SortOrder: 1,
			Enabled:   true,
		},
	}
}

func newRecommendationBody(id, metric string) RecommendationCreate {
	when := GroupWhen{Match: MatchAll, Predicates: []Predicate{}}
	if metric != "" {
		when.Predicates = append(when.Predicates, Predicate{Metric: metric, Op: OpGt, Value: 0})
	}

	return RecommendationCreate{
		ID: id,
		RecommendationWrite: RecommendationWrite{
			Title:     "Продвиньте объявления",
			Text:      "Продвижение поможет найти покупателя быстрее",
			Callout:   "Больше показов",
			LinkLabel: "Продвинуть",
			Path:      "/profile",
			Priority:  100,
			When:      when,
			Enabled:   true,
		},
	}
}

func newMetricBody(key string) MetricCreate {
	return MetricCreate{
		Key: key,
		MetricWrite: MetricWrite{
			ValueType:    MetricValueNumber,
			IsPublic:     true,
			IncludeInLlm: true,
			SortOrder:    10,
			Enabled:      true,
		},
	}
}
