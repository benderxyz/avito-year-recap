package cards

import (
	"strings"
	"testing"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

func sample(value float64) clients.MetricSample {
	return clients.MetricSample{Value: &value}
}

func sampleWithPercentile(value, percentile float64) clients.MetricSample {
	return clients.MetricSample{Value: &value, Percentile: &percentile}
}

func fullMetrics() clients.Metrics {
	return clients.Metrics{
		"listingsPublished": sampleWithPercentile(12, 88),
		"viewsTotal":        sampleWithPercentile(3400, 92),
		"favoritesReceived": sampleWithPercentile(87, 79),
		"messagesSent":      sampleWithPercentile(1250, 85),
		"dealsClosed":       sampleWithPercentile(19, 74),
		"moneyEarned":       sample(150000),
		"moneySaved":        sample(32000),
		"daysActive":        sample(210),
		"peakDayViews":      sample(340),
		"searchQueries":     sample(500),
		"categoriesTried":   sample(7),
		"deliveryOrders":    sample(5),
		"activeListings":    sample(8),
		"avgReplySeconds":   sample(120),
		"firstListingAt":    sample(1704067200),
		"firstDealAt":       sample(1706745600),
	}
}

func privateOptions(key []byte) BuildOptions {
	rules := testRuleSet()
	return BuildOptions{
		Mode:         models.RecapModePrivate,
		SigningKey:   key,
		ShareBaseURL: "http://localhost:3000",
		Rules:        &rules,
	}
}

func publicOptions() BuildOptions {
	rules := testRuleSet()
	return BuildOptions{
		Mode:  models.RecapModePublic,
		Rules: &rules,
	}
}

func privateOptionsWithThreshold(metricKey string, threshold float64) BuildOptions {
	rules := testRuleSet()
	for i := range rules.metrics {
		if rules.metrics[i].Key == metricKey {
			rules.metrics[i].ComparisonMinPercentile = threshold
		}
	}
	return BuildOptions{
		Mode:         models.RecapModePrivate,
		SigningKey:   []byte("k"),
		ShareBaseURL: "http://localhost:3000",
		Rules:        &rules,
	}
}

func listingsScene(t *testing.T, story []map[string]any) map[string]any {
	t.Helper()
	for _, scene := range story {
		if scene["id"] == "stat-listings" {
			return scene
		}
	}
	t.Fatal("expected stat-listings scene")
	return nil
}

func storyIDs(story []map[string]any) []string {
	ids := make([]string, 0, len(story))
	for _, scene := range story {
		if id, ok := scene["id"].(string); ok {
			ids = append(ids, id)
		}
	}
	return ids
}

func TestBuildRecapPrivateShouldExposeSensitiveMetrics(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	for _, name := range []string{
		"moneyEarned", "moneySaved", "messagesSent",
		"avgReplySeconds", "activeListings",
		"firstListingAt", "firstDealAt",
	} {
		if _, ok := recap.Metrics[name]; !ok {
			t.Fatalf("expected %s in private metrics", name)
		}
	}
}

func TestBuildRecapPrivateShouldIncludePercentileMetrics(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	for _, name := range []string{"listingsPercentile", "viewsPercentile", "messagesPercentile"} {
		if _, ok := recap.Metrics[name]; !ok {
			t.Fatalf("expected %s in private metrics", name)
		}
	}
}

func TestBuildRecapPrivateShouldAttachComparisonToStatScenes(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	scene := listingsScene(t, recap.Story)
	if scene["percentile"] != "listingsPercentile" {
		t.Fatalf("expected listings percentile key in story scene, got %v", scene["percentile"])
	}
	if scene["comparisonTemplate"] == nil {
		t.Fatal("expected comparison template in listings scene")
	}
}

func TestBuildRecapShouldSkipZeroPercentileComparison(t *testing.T) {
	metrics := fullMetrics()
	metrics["listingsPublished"] = sampleWithPercentile(12, 0)
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, metrics, privateOptions([]byte("k")))

	for _, scene := range recap.Story {
		if scene["id"] == "stat-listings" {
			if _, ok := scene["percentile"]; ok {
				t.Fatal("did not expect zero percentile comparison")
			}
			return
		}
	}
	t.Fatal("expected stat-listings scene")
}

func TestBuildRecapShouldSkipComparisonWhenPercentileBelowMetricThreshold(t *testing.T) {
	metrics := fullMetrics()
	metrics["listingsPublished"] = sampleWithPercentile(12, 88)
	opts := privateOptionsWithThreshold("listingsPublished", 90)

	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, metrics, opts)

	if _, ok := listingsScene(t, recap.Story)["percentile"]; ok {
		t.Fatal("did not expect comparison below metric threshold")
	}
}

func TestBuildRecapShouldAttachComparisonWhenPercentileEqualsMetricThreshold(t *testing.T) {
	metrics := fullMetrics()
	metrics["listingsPublished"] = sampleWithPercentile(12, 88)
	opts := privateOptionsWithThreshold("listingsPublished", 88)

	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, metrics, opts)

	if listingsScene(t, recap.Story)["percentile"] != "listingsPercentile" {
		t.Fatal("expected comparison at metric threshold boundary")
	}
}

func TestBuildRecapShouldSkipZeroPercentileEvenWhenThresholdIsZero(t *testing.T) {
	metrics := fullMetrics()
	metrics["listingsPublished"] = sampleWithPercentile(12, 0)
	opts := privateOptionsWithThreshold("listingsPublished", 0)

	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, metrics, opts)

	if _, ok := listingsScene(t, recap.Story)["percentile"]; ok {
		t.Fatal("did not expect comparison for zero percentile")
	}
}

func TestBuildRecapPublicShouldFilterSensitiveMetrics(t *testing.T) {
	rules := testRuleSet()
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), publicOptions())

	for _, name := range []string{"moneyEarned", "moneySaved", "messagesSent", "avgReplySeconds", "activeListings"} {
		if _, ok := recap.Metrics[name]; ok {
			t.Fatalf("did not expect %s in public metrics", name)
		}
	}

	publicKeys := make(map[string]bool, len(rules.metrics))
	for _, def := range rules.metrics {
		if def.IsPublic {
			publicKeys[def.Key] = true
		}
	}

	for name := range recap.Metrics {
		if !publicKeys[name] {
			t.Fatalf("metric %q is not marked public in metric definitions", name)
		}
	}
}

func TestBuildRecapPublicShouldOmitPrivateScenes(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), publicOptions())

	for _, scene := range recap.Story {
		id, _ := scene["id"].(string)
		switch id {
		case "stat-messages", "stat-earned", "stat-saved", "stat-reply", "insight-first-listing", "insight-first-deal":
			t.Fatalf("private scene %v leaked into public story", id)
		}
	}
}

func TestBuildRecapPrivateShouldIncludeShareFeatures(t *testing.T) {
	recap := BuildRecap(models.Profile{UserID: 42, ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	if recap.Features == nil || !recap.Features.ShareEnabled {
		t.Fatal("expected share features in private recap")
	}
	if !strings.HasPrefix(recap.Features.ShareURL, "http://localhost:3000/share/") {
		t.Fatalf("unexpected share url %q", recap.Features.ShareURL)
	}

	id, year, err := DecodeShareToken([]byte("k"), strings.TrimPrefix(recap.Features.ShareURL, "http://localhost:3000/share/"))
	if err != nil {
		t.Fatalf("share url token should be valid: %v", err)
	}
	if id != "42" {
		t.Fatalf("expected numeric user id 42 in share token, got %q", id)
	}
	if year != 2024 {
		t.Fatalf("expected year 2024 in share token, got %d", year)
	}
}

func TestBuildRecapPublicShouldNotExposeShareLink(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), publicOptions())

	if recap.Features != nil && recap.Features.ShareEnabled {
		t.Fatal("public recap must not advertise sharing")
	}
}

func TestBuildRecapBadgeShouldMaskExactCountInPublic(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), publicOptions())

	for _, badge := range recap.Badges {
		if badge.ID == "messages_sent_1000" {
			if strings.Contains(badge.Description, "1250") {
				t.Fatalf("public badge leaked exact count: %q", badge.Description)
			}
		}
	}
}

func hasRecommendationScene(story []map[string]any) bool {
	for _, scene := range story {
		if id, ok := scene["id"].(string); ok && strings.HasPrefix(id, "reco-") {
			return true
		}
	}
	return false
}

func TestBuildRecapPrivateShouldIncludeRecommendations(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	if !hasRecommendationScene(recap.Story) {
		t.Fatal("expected recommendation scenes in private recap")
	}
}

func TestBuildRecapPublicShouldOmitRecommendations(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), publicOptions())

	if hasRecommendationScene(recap.Story) {
		t.Fatal("public share-card must not include recommendation scenes")
	}
}

func TestBuildRecapPrivateShouldRenderIntroFromDictionary(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	if len(recap.Story) == 0 || recap.Story[0]["id"] != "intro" {
		t.Fatalf("expected intro first, got %v", storyIDs(recap.Story))
	}
	if recap.Story[0]["title"] != "alex, ваш 2024 на Авито" {
		t.Fatalf("unexpected intro title %v", recap.Story[0]["title"])
	}
}

func TestBuildRecapPrivateShouldIncludeScenesForCalculatedMetrics(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))
	ids := storyIDs(recap.Story)

	for _, want := range []string{
		"stat-saved", "blocks-days-active", "stat-peak-views", "stat-categories",
		"blocks-search", "stat-delivery", "stat-active-listings",
		"stat-reply", "insight-first-listing", "insight-first-deal",
	} {
		found := false
		for _, id := range ids {
			if id == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected scene %q in story, got %v", want, ids)
		}
	}
}

func findScene(story []map[string]any, id string) map[string]any {
	for _, scene := range story {
		if scene["id"] == id {
			return scene
		}
	}
	return nil
}

func TestBuildRecapShouldLeaveValueTemplateToFrontendWhenSceneUsesMetric(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	scene := findScene(recap.Story, "insight-first-listing")
	if scene == nil {
		t.Fatal("expected insight-first-listing scene")
	}
	if scene["text"] != "Первое объявление в этом году вы опубликовали {{value}}" {
		t.Fatalf("expected untouched value template, got %v", scene["text"])
	}
	if scene["value"] != "firstListingAt" {
		t.Fatalf("expected metric key in scene value, got %v", scene["value"])
	}
}

func TestBuildRecapShouldEmitDateMetricAsISOStringWhenDefinitionIsDate(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	listing, ok := recap.Metrics["firstListingAt"]
	if !ok {
		t.Fatal("expected firstListingAt metric")
	}
	if listing.Type != models.MetricTypeDate {
		t.Fatalf("expected date metric type, got %q", listing.Type)
	}
	if listing.Value != "2024-01-01" {
		t.Fatalf("unexpected firstListingAt value %v", listing.Value)
	}
}

func TestBuildRecapShouldSkipMetricWhenAnalyticsOmitsSourceKey(t *testing.T) {
	metrics := fullMetrics()
	delete(metrics, "firstDealAt")

	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, metrics, privateOptions([]byte("k")))

	if _, ok := recap.Metrics["firstDealAt"]; ok {
		t.Fatal("did not expect firstDealAt metric when analytics omits it")
	}
	if findScene(recap.Story, "insight-first-deal") != nil {
		t.Fatal("did not expect insight-first-deal scene without the metric")
	}
}

func TestBuildRecapShouldExposeMetricWhenDefinitionIsAddedWithoutCodeChanges(t *testing.T) {
	rules := testRuleSet()
	rules.metrics = append(rules.metrics, models.MetricDefinition{
		Key:         "storiesWatched",
		ValueType:   models.MetricTypeNumber,
		IsPublic:    true,
		SourceKey:   "storiesWatched",
		SourceField: models.MetricSourceValue,
	})

	metrics := fullMetrics()
	metrics["storiesWatched"] = sample(9)

	recap := BuildRecap(
		models.Profile{ExternalID: "u1", Username: "alex"},
		2024,
		metrics,
		BuildOptions{Mode: models.RecapModePrivate, Rules: &rules},
	)

	metric, ok := recap.Metrics["storiesWatched"]
	if !ok {
		t.Fatal("expected metric defined only in metric definitions")
	}
	if metric.Value != float64(9) {
		t.Fatalf("unexpected storiesWatched value %v", metric.Value)
	}
}

func TestBuildRecapShouldKeepRuleComparisonTemplateInsteadOfDefault(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	want := "Объявлений у вас больше, чем у {{percentile}}% продавцов"
	if got := listingsScene(t, recap.Story)["comparisonTemplate"]; got != want {
		t.Fatalf("want %q, got %v", want, got)
	}
}

func TestBuildRecapShouldUseDefaultComparisonTemplateWhenRuleHasNone(t *testing.T) {
	rules := testRuleSet()
	for i := range rules.stories {
		delete(rules.stories[i].scene, "comparisonTemplate")
	}
	opts := BuildOptions{Mode: models.RecapModePrivate, SigningKey: []byte("k"), Rules: &rules}

	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), opts)

	want := "Больше, чем у {{percentile}}% пользователей Авито"
	if got := listingsScene(t, recap.Story)["comparisonTemplate"]; got != want {
		t.Fatalf("want %q, got %v", want, got)
	}
}

func TestBuildRecapShouldRenderCategoriesAsBlocksWithCallout(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	var scene map[string]any
	for _, s := range recap.Story {
		if s["id"] == "stat-categories" {
			scene = s
			break
		}
	}
	if scene == nil {
		t.Fatal("expected stat-categories scene")
	}
	if scene["type"] != "blocks" {
		t.Fatalf("want blocks scene, got %v", scene["type"])
	}

	blocks, ok := scene["blocks"].([]any)
	if !ok || len(blocks) != 2 {
		t.Fatalf("want two blocks, got %v", scene["blocks"])
	}

	stat, _ := blocks[0].(map[string]any)
	if stat["title"] != "вы открывали за год" {
		t.Fatalf("unexpected stat title %v", stat["title"])
	}

	callout, _ := blocks[1].(map[string]any)
	if callout["type"] != "callout" || callout["text"] == "" {
		t.Fatalf("want non-empty callout block, got %v", blocks[1])
	}
}
