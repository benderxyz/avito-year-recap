package cards

import (
	"strings"
	"testing"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

func fullMetrics() clients.Metrics {
	return clients.Metrics{
		ListingsPublished: 12,
		ViewsTotal:        3400,
		FavoritesReceived: 87,
		MessagesSent:      1250,
		DealsClosed:       19,
		MoneyEarned:       150000,
		MoneySaved:        32000,
		DaysActive:        210,
		PeakDayViews:      340,
		SearchQueries:     500,
		CategoriesTried:   7,
		DeliveryOrders:    5,
		ActiveListings:    8,
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
