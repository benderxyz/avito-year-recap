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
	}
}

func privateOptions(key []byte) BuildOptions {
	return BuildOptions{
		Mode:         models.RecapModePrivate,
		SigningKey:   key,
		ShareBaseURL: "http://localhost:3000",
	}
}

func TestBuildRecapPrivateShouldExposeSensitiveMetrics(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	for _, name := range []string{"moneyEarned", "moneySaved", "messagesSent"} {
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
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), BuildOptions{Mode: models.RecapModePublic})

	for _, name := range []string{"moneyEarned", "moneySaved", "messagesSent", "avgReplySeconds"} {
		if _, ok := recap.Metrics[name]; ok {
			t.Fatalf("did not expect %s in public metrics", name)
		}
	}

	for name := range recap.Metrics {
		if !models.PublicMetricsAllowlist[name] {
			t.Fatalf("metric %q is not in the public allowlist", name)
		}
	}
}

func TestBuildRecapPublicShouldOmitPrivateScenes(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), BuildOptions{Mode: models.RecapModePublic})

	for _, scene := range recap.Story {
		if scene["id"] == "stat-messages" || scene["id"] == "stat-earned" {
			t.Fatalf("private scene %v leaked into public story", scene["id"])
		}
	}
}

func TestBuildRecapPrivateShouldIncludeShareFeatures(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), privateOptions([]byte("k")))

	if recap.Features == nil || !recap.Features.ShareEnabled {
		t.Fatal("expected share features in private recap")
	}
	if !strings.HasPrefix(recap.Features.ShareURL, "http://localhost:3000/share/") {
		t.Fatalf("unexpected share url %q", recap.Features.ShareURL)
	}
	if _, _, err := DecodeShareToken([]byte("k"), strings.TrimPrefix(recap.Features.ShareURL, "http://localhost:3000/share/")); err != nil {
		t.Fatalf("share url token should be valid: %v", err)
	}
}

func TestBuildRecapPublicShouldNotExposeShareLink(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), BuildOptions{Mode: models.RecapModePublic})

	if recap.Features != nil && recap.Features.ShareEnabled {
		t.Fatal("public recap must not advertise sharing")
	}
}

func TestBuildRecapBadgeShouldMaskExactCountInPublic(t *testing.T) {
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), BuildOptions{Mode: models.RecapModePublic})

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
	recap := BuildRecap(models.Profile{ExternalID: "u1", Username: "alex"}, 2024, fullMetrics(), BuildOptions{Mode: models.RecapModePublic})

	if hasRecommendationScene(recap.Story) {
		t.Fatal("public share-card must not include recommendation scenes")
	}
}
