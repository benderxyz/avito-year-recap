package cards

import (
	"strings"
	"testing"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

func fullMetrics() clients.Metrics {
	listingsPercentile := float64(88)
	viewsPercentile := float64(92)
	favoritesPercentile := float64(79)
	messagesPercentile := float64(85)
	dealsPercentile := float64(74)

	return clients.Metrics{
		ListingsPublished:   12,
		ListingsPercentile:  &listingsPercentile,
		ViewsTotal:          3400,
		ViewsPercentile:     &viewsPercentile,
		FavoritesReceived:   87,
		FavoritesPercentile: &favoritesPercentile,
		MessagesSent:        1250,
		MessagesPercentile:  &messagesPercentile,
		DealsClosed:         19,
		DealsPercentile:     &dealsPercentile,
		MoneyEarned:         150000,
		MoneySaved:          32000,
		DaysActive:          210,
		PeakDayViews:        340,
		SearchQueries:       500,
		CategoriesTried:     7,
		DeliveryOrders:      5,
		ActiveListings:      8,
		SellerRating:        4.8,
		AvgReplySeconds:     120,
		FirstListingAt:      1704067200,
		FirstDealAt:         1706745600,
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
		"sellerRating", "avgReplySeconds", "activeListings",
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

	for _, scene := range recap.Story {
		if scene["id"] == "stat-listings" {
			if scene["percentile"] != "listingsPercentile" {
				t.Fatalf("expected listings percentile key in story scene, got %v", scene["percentile"])
			}
			if scene["comparisonTemplate"] == nil {
				t.Fatal("expected comparison template in listings scene")
			}
		}
	}
}

func TestBuildRecapPublicShouldFilterSensitiveMetrics(t *testing.T) {
	rules := testRuleSet()
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), publicOptions())

	for _, name := range []string{"moneyEarned", "moneySaved", "messagesSent", "avgReplySeconds", "sellerRating", "activeListings"} {
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
		case "stat-messages", "stat-earned", "stat-saved", "stat-rating", "stat-reply", "insight-first-listing", "insight-first-deal":
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
		"blocks-search", "stat-delivery", "stat-active-listings", "stat-rating",
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

func TestBuildRecapPrivateShouldFormatInsightDates(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	for _, scene := range recap.Story {
		if scene["id"] == "insight-first-listing" {
			text, _ := scene["text"].(string)
			if !strings.Contains(text, "2024-01-01") {
				t.Fatalf("expected formatted first listing date in insight text, got %q", text)
			}
		}
	}
}

func TestBuildRecapPrivateShouldEmitStringMilestoneMetrics(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	listing, ok := recap.Metrics["firstListingAt"]
	if !ok {
		t.Fatal("expected firstListingAt metric")
	}
	if listing.Type != models.MetricTypeString {
		t.Fatalf("expected string metric type, got %q", listing.Type)
	}
	if listing.Value != "2024-01-01" {
		t.Fatalf("unexpected firstListingAt value %v", listing.Value)
	}
}
