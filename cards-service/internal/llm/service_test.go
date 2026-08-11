package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

type fakeProvider struct {
	calls    int
	response string
	err      error
	block    bool
}

func (f *fakeProvider) Name() string { return "fake" }

func (f *fakeProvider) Complete(ctx context.Context, _ CompletionRequest) (string, error) {
	f.calls++
	if f.block {
		<-ctx.Done()
		return "", ctx.Err()
	}
	if f.err != nil {
		return "", f.err
	}
	return f.response, nil
}

type fakeCache struct {
	data     map[string]Result
	getCalls int
	putCalls int
}

func newFakeCache() *fakeCache { return &fakeCache{data: map[string]Result{}} }

func cacheKey(userID string, year int, mode, version string) string {
	return fmt.Sprintf("%s:%d:%s:%s", userID, year, mode, version)
}

func (c *fakeCache) Get(_ context.Context, userID string, year int, mode, version string) (Result, bool, error) {
	c.getCalls++
	r, ok := c.data[cacheKey(userID, year, mode, version)]
	return r, ok, nil
}

func (c *fakeCache) Put(_ context.Context, userID string, year int, mode, version string, result Result) error {
	c.putCalls++
	c.data[cacheKey(userID, year, mode, version)] = result
	return nil
}

func basePayload() models.RecapPayload {
	return models.RecapPayload{
		Meta: models.Meta{
			Year: 2024,
			User: models.User{ID: "u1", DisplayName: "Аня"},
		},
		Badges: []models.Badge{
			{ID: "seller", Title: "Продавец", Description: "Закрыл много сделок"},
			{ID: "explorer", Title: "Исследователь", Description: "Пробовал разные категории"},
		},
		Story: []map[string]any{
			{"id": "intro", "type": "intro", "title": "Привет"},
			{"id": "stat-1", "type": "stat", "label": "объявлений"},
			{"id": "outro", "type": "outro"},
		},
	}
}

const validResponse = `{
  "badges": {
    "seller": {"title": "Мастер сделок", "description": "Каждая продажа как маленькая победа"},
    "explorer": {"title": "Вечный искатель", "description": "Тебе интересно всё вокруг"}
  },
  "insight": {"title": "Вдумчивый охотник", "text": "Ты редко пишешь, но точно знаешь, что ищешь"}
}`

func newTestService(p Provider, c cacheStore) *Service {
	return &Service{provider: p, cache: c, timeout: time.Second, version: "testver"}
}

func sceneIndexByID(story []map[string]any, id string) int {
	for i, s := range story {
		if sid, _ := s["id"].(string); sid == id {
			return i
		}
	}
	return -1
}

func badgeByID(badges []models.Badge, id string) (models.Badge, bool) {
	for _, b := range badges {
		if b.ID == id {
			return b, true
		}
	}
	return models.Badge{}, false
}

func TestEnrichValidRewritesBadgesAndInsertsInsight(t *testing.T) {
	for _, mode := range []models.RecapMode{models.RecapModePrivate, models.RecapModePublic} {
		t.Run(string(mode), func(t *testing.T) {
			provider := &fakeProvider{response: validResponse}
			svc := newTestService(provider, newFakeCache())

			out, report, err := svc.Enrich(context.Background(),
				EnrichInput{ID: "u1", Year: 2024, Mode: mode}, basePayload())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if report.Source != "llm" {
				t.Errorf("source = %q, want llm", report.Source)
			}
			if report.BadgesRewritten != 2 || report.BadgesFallback != 0 {
				t.Errorf("rewritten=%d fallback=%d, want 2/0", report.BadgesRewritten, report.BadgesFallback)
			}
			if !report.InsightAdded {
				t.Error("InsightAdded = false, want true")
			}

			seller, _ := badgeByID(out.Badges, "seller")
			if seller.Title != "Мастер сделок" {
				t.Errorf("seller title = %q, want rewritten", seller.Title)
			}

			idx := sceneIndexByID(out.Story, "llm-insight")
			if idx != 1 {
				t.Fatalf("llm-insight at index %d, want 1 (after intro)", idx)
			}
			if out.Story[idx]["type"] != "insight" {
				t.Errorf("insight scene type = %v, want insight", out.Story[idx]["type"])
			}
		})
	}
}

func TestEnrichInvalidBadgeFallsBackToDefault(t *testing.T) {
	resp := `{
	  "badges": {
	    "seller": {"title": "Топ 1 продавец", "description": "Лучший результат"},
	    "explorer": {"title": "Вечный искатель", "description": "Тебе интересно всё вокруг"}
	  },
	  "insight": null
	}`
	provider := &fakeProvider{response: resp}
	svc := newTestService(provider, newFakeCache())

	out, report, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024}, basePayload())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.BadgesRewritten != 1 || report.BadgesFallback != 1 {
		t.Errorf("rewritten=%d fallback=%d, want 1/1", report.BadgesRewritten, report.BadgesFallback)
	}

	seller, _ := badgeByID(out.Badges, "seller")
	if seller.Title != "Продавец" {
		t.Errorf("seller title = %q, want default kept", seller.Title)
	}
	explorer, _ := badgeByID(out.Badges, "explorer")
	if explorer.Title != "Вечный искатель" {
		t.Errorf("explorer title = %q, want rewritten", explorer.Title)
	}
	if report.InsightAdded {
		t.Error("InsightAdded = true, want false for null insight")
	}
}

func TestEnrichProviderErrorReturnsOriginalPayloadWithError(t *testing.T) {
	provider := &fakeProvider{err: errors.New("boom")}
	svc := newTestService(provider, newFakeCache())

	in := basePayload()
	out, report, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024}, in)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if report.Source != "llm" {
		t.Errorf("source = %q, want llm", report.Source)
	}
	seller, _ := badgeByID(out.Badges, "seller")
	if seller.Title != "Продавец" {
		t.Errorf("seller title = %q, want default", seller.Title)
	}
	if sceneIndexByID(out.Story, "llm-insight") != -1 {
		t.Error("insight scene present, want none on failure")
	}
	if len(out.Story) != 3 {
		t.Errorf("story len = %d, want 3 (unchanged)", len(out.Story))
	}
}

func TestEnrichTimeoutReturnsError(t *testing.T) {
	provider := &fakeProvider{block: true}
	svc := &Service{provider: provider, cache: newFakeCache(), timeout: 20 * time.Millisecond}

	_, report, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024}, basePayload())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if report.Source != "llm" {
		t.Errorf("source = %q, want llm", report.Source)
	}
	if provider.calls != 1 {
		t.Errorf("provider called %d times, want 1", provider.calls)
	}
}

func TestEnrichSecondCallHitsCacheAndSkipsProvider(t *testing.T) {
	provider := &fakeProvider{response: validResponse}
	cache := newFakeCache()
	svc := newTestService(provider, cache)

	if _, _, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024}, basePayload()); err != nil {
		t.Fatalf("first enrich: %v", err)
	}
	if provider.calls != 1 || cache.putCalls != 1 {
		t.Fatalf("after first: calls=%d put=%d, want 1/1", provider.calls, cache.putCalls)
	}

	out, report, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024}, basePayload())
	if err != nil {
		t.Fatalf("second enrich: %v", err)
	}
	if provider.calls != 1 {
		t.Errorf("provider called %d times, want 1 (cache hit)", provider.calls)
	}
	if report.Source != "cache" {
		t.Errorf("source = %q, want cache", report.Source)
	}
	seller, _ := badgeByID(out.Badges, "seller")
	if seller.Title != "Мастер сделок" {
		t.Errorf("seller title = %q, want rewritten from cache", seller.Title)
	}
	if sceneIndexByID(out.Story, "llm-insight") != 1 {
		t.Error("insight scene missing on cache hit")
	}
}

func TestEnrichModeKeysCacheSeparately(t *testing.T) {
	provider := &fakeProvider{response: validResponse}
	cache := newFakeCache()
	svc := newTestService(provider, cache)

	if _, _, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024, Mode: models.RecapModePrivate}, basePayload()); err != nil {
		t.Fatalf("private enrich: %v", err)
	}
	_, report, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024, Mode: models.RecapModePublic}, basePayload())
	if err != nil {
		t.Fatalf("public enrich: %v", err)
	}
	if report.Source != "llm" {
		t.Errorf("public source = %q, want llm (private entry must not leak to public)", report.Source)
	}
	if provider.calls != 2 {
		t.Errorf("provider calls = %d, want 2 (one per mode)", provider.calls)
	}
	if len(cache.data) != 2 {
		t.Errorf("cache entries = %d, want 2 (separate per mode)", len(cache.data))
	}
}

func TestEnrichPromptVersionInvalidatesCache(t *testing.T) {
	provider := &fakeProvider{response: validResponse}
	cache := newFakeCache()
	svc := newTestService(provider, cache)

	if _, _, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024, Mode: models.RecapModePrivate}, basePayload()); err != nil {
		t.Fatalf("first enrich: %v", err)
	}
	svc.version = "newver"
	_, report, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024, Mode: models.RecapModePrivate}, basePayload())
	if err != nil {
		t.Fatalf("second enrich: %v", err)
	}
	if report.Source != "llm" {
		t.Errorf("source = %q, want llm (version change must invalidate cache)", report.Source)
	}
	if provider.calls != 2 {
		t.Errorf("provider calls = %d, want 2 (version change → regenerate)", provider.calls)
	}
}

func TestPromptVersionStableAndModelSensitive(t *testing.T) {
	a := promptVersion("model-a")
	if len(a) != 12 {
		t.Errorf("promptVersion len = %d, want 12", len(a))
	}
	if a != promptVersion("model-a") {
		t.Error("promptVersion should be stable for the same model")
	}
	if a == promptVersion("model-b") {
		t.Error("promptVersion should differ across models")
	}
}

func TestSanitizeIgnoresHallucinatedBadgeIDs(t *testing.T) {
	resp := llmResponse{
		Badges: map[string]BadgeText{
			"seller": {Title: "Мастер сделок", Description: "Каждая продажа как победа"},
			"ghost":  {Title: "Призрак", Description: "Такого бейджа нет"},
		},
	}
	result := sanitize(resp, basePayload().Badges)
	if _, ok := result.Badges["ghost"]; ok {
		t.Error("hallucinated id 'ghost' leaked into result")
	}
	if _, ok := result.Badges["seller"]; !ok {
		t.Error("valid id 'seller' missing from result")
	}
}

func TestEnrichDropsInsightLeakingSpecifics(t *testing.T) {
	resp := `{
	  "badges": {},
	  "insight": {"title": "Итог", "text": "Ты заработал 100000 за год"}
	}`
	provider := &fakeProvider{response: resp}
	svc := newTestService(provider, newFakeCache())

	out, report, err := svc.Enrich(context.Background(),
		EnrichInput{ID: "u1", Year: 2024, Mode: models.RecapModePublic}, basePayload())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.InsightAdded {
		t.Error("InsightAdded = true, want false for leaking insight")
	}
	if sceneIndexByID(out.Story, "llm-insight") != -1 {
		t.Error("leaking insight scene present, want dropped")
	}
}

func TestNilServiceAndNilProviderAreNoOps(t *testing.T) {
	var nilSvc *Service
	in := basePayload()
	out, report, err := nilSvc.Enrich(context.Background(), EnrichInput{ID: "u1", Year: 2024}, in)
	if err != nil || report.Source != "skip" {
		t.Errorf("nil service: err=%v source=%q, want nil/skip", err, report.Source)
	}
	if len(out.Story) != 3 {
		t.Errorf("nil service mutated story: len=%d", len(out.Story))
	}

	svc := newTestService(nil, newFakeCache())
	_, report, err = svc.Enrich(context.Background(), EnrichInput{ID: "u1", Year: 2024}, basePayload())
	if err != nil || report.Source != "skip" {
		t.Errorf("nil provider: err=%v source=%q, want nil/skip", err, report.Source)
	}
}

func TestInsertInsightWithoutIntroGoesToFront(t *testing.T) {
	story := []map[string]any{
		{"id": "stat-1", "type": "stat"},
		{"id": "outro", "type": "outro"},
	}
	out := insertInsight(story, InsightText{Title: "T", Text: "X"})
	if out[0]["id"] != "llm-insight" {
		t.Errorf("first scene id = %v, want llm-insight", out[0]["id"])
	}
	if len(out) != 3 {
		t.Errorf("len = %d, want 3", len(out))
	}
}

func promptSample(value float64) clients.MetricSample {
	return clients.MetricSample{Value: &value}
}

func promptDefinition(key string, isPublic, includeInLLM bool) models.MetricDefinition {
	return models.MetricDefinition{
		Key:          key,
		ValueType:    models.MetricTypeNumber,
		IsPublic:     isPublic,
		SourceKey:    key,
		SourceField:  models.MetricSourceValue,
		IncludeInLLM: includeInLLM,
	}
}

func TestMetricsSummaryShouldKeepOnlyPublicMetricsMarkedForPrompt(t *testing.T) {
	m := clients.Metrics{
		"listingsPublished": promptSample(5),
		"viewsTotal":        promptSample(0),
		"messagesSent":      promptSample(3),
		"sellerRating":      promptSample(4.8),
		"avgReplySeconds":   promptSample(120),
	}
	defs := []models.MetricDefinition{
		promptDefinition("listingsPublished", true, true),
		promptDefinition("viewsTotal", true, true),
		promptDefinition("messagesSent", false, true),
		promptDefinition("sellerRating", false, true),
		promptDefinition("avgReplySeconds", true, false),
	}

	joined := strings.Join(metricsSummary(m, defs), "\n")

	if !strings.Contains(joined, "listingsPublished: 5") {
		t.Errorf("missing public listingsPublished; got %q", joined)
	}
	if strings.Contains(joined, "viewsTotal") {
		t.Errorf("zero-valued public metric leaked; got %q", joined)
	}
	if strings.Contains(joined, "sellerRating") {
		t.Errorf("private sellerRating leaked into prompt; got %q", joined)
	}
	if strings.Contains(joined, "messagesSent") {
		t.Errorf("private messagesSent leaked into prompt; got %q", joined)
	}
	if strings.Contains(joined, "avgReplySeconds") {
		t.Errorf("metric not marked for prompt leaked; got %q", joined)
	}
}

func TestMetricsSummaryShouldFailClosedWhenDefinitionsAreEmpty(t *testing.T) {
	m := clients.Metrics{"listingsPublished": promptSample(5), "sellerRating": promptSample(4.8)}

	if lines := metricsSummary(m, nil); len(lines) != 0 {
		t.Errorf("nil definitions should yield no lines, got %v", lines)
	}
	if lines := metricsSummary(m, []models.MetricDefinition{}); len(lines) != 0 {
		t.Errorf("empty definitions should yield no lines, got %v", lines)
	}
}

func TestMetricsSummaryShouldKeepOneDecimalWhenValueIsFractional(t *testing.T) {
	m := clients.Metrics{"sellerRating": promptSample(4.8)}
	defs := []models.MetricDefinition{promptDefinition("sellerRating", true, true)}

	lines := metricsSummary(m, defs)

	if len(lines) != 1 || lines[0] != "sellerRating: 4.8" {
		t.Errorf("expected fractional formatting, got %v", lines)
	}
}

func TestExtractJSONObject(t *testing.T) {
	in := "Вот ответ: {\"badges\":{}} — готово"
	if got := extractJSONObject(in); got != `{"badges":{}}` {
		t.Errorf("extractJSONObject = %q", got)
	}
}

func TestStripCodeFence(t *testing.T) {
	in := "```json\n{\"a\":1}\n```"
	if got := stripCodeFence(in); got != `{"a":1}` {
		t.Errorf("stripCodeFence = %q, want {\"a\":1}", got)
	}
}
