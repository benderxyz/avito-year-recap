package cards

import (
	"strings"
	"testing"

	"cards-service/internal/clients"
)

func recoIDs(scenes []map[string]any) []string {
	ids := make([]string, 0, len(scenes))
	for _, scene := range scenes {
		if id, ok := scene["id"].(string); ok {
			ids = append(ids, id)
		}
	}
	return ids
}

func TestBuildRecommendationsShouldCapAndRankByPriority(t *testing.T) {
	scenes := buildRecommendations(defaultRecommendationRules, fullMetrics(), "https://www.avito.ru")

	if len(scenes) == 0 {
		t.Fatal("expected recommendations for an active user")
	}
	if len(scenes) > maxRecommendations {
		t.Fatalf("expected at most %d recommendations, got %d", maxRecommendations, len(scenes))
	}

	ids := recoIDs(scenes)
	if ids[0] != "reco-listings-favorited" {
		t.Fatalf("expected highest-priority matched rule first, got %q", ids[0])
	}
}

func TestBuildRecommendationsShouldSuggestPromotionForSellerWithoutDeals(t *testing.T) {
	m := clients.Metrics{ListingsPublished: 5, DealsClosed: 0}

	ids := recoIDs(buildRecommendations(defaultRecommendationRules, m, "https://www.avito.ru"))

	if ids[0] != "reco-resume-listings" {
		t.Fatalf("expected promotion recommendation first for seller without deals, got %v", ids)
	}
}

func TestBuildRecommendationsShouldSuggestFirstListingForInactiveUser(t *testing.T) {
	ids := recoIDs(buildRecommendations(defaultRecommendationRules, clients.Metrics{}, "https://www.avito.ru"))

	found := false
	for _, id := range ids {
		if id == "reco-new-listing" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected first-listing recommendation for inactive user, got %v", ids)
	}
}

func TestBuildRecommendationsShouldBuildAbsoluteAvitoLinks(t *testing.T) {
	scenes := buildRecommendations(defaultRecommendationRules, clients.Metrics{FavoritesReceived: 3}, "https://www.avito.ru")

	actions, ok := scenes[0]["actions"].([]map[string]any)
	if !ok || len(actions) == 0 {
		t.Fatal("expected recommendation actions")
	}

	href, _ := actions[0]["href"].(string)
	if !strings.HasPrefix(href, "https://www.avito.ru/") {
		t.Fatalf("expected absolute avito href, got %q", href)
	}
}

func TestBuildProductURLShouldNotDoubleSlash(t *testing.T) {
	got := buildProductURL("https://www.avito.ru/", "/favorites")
	if got != "https://www.avito.ru/favorites" {
		t.Fatalf("unexpected url %q", got)
	}
}
