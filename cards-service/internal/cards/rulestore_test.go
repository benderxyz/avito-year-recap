package cards

import (
	"context"
	"testing"
	"time"

	"cards-service/internal/models"
)

func TestMakeBadgeRuleShouldBuildSinglePredicateCondition(t *testing.T) {
	rule := makeBadgeRule("messages_sent_1000", "Разговорчивый", "text", "/badges/messages_sent_1000.svg", "both", "messagesSent", "gt", 1000)

	if rule.badge.ID != "messages_sent_1000" || rule.badge.Title != "Разговорчивый" {
		t.Fatalf("unexpected badge mapping: %+v", rule.badge)
	}
	if rule.badge.IconURL != "/badges/messages_sent_1000.svg" {
		t.Fatalf("unexpected badge icon %q", rule.badge.IconURL)
	}
	if rule.visibility != visibilityBoth {
		t.Fatalf("unexpected visibility %q", rule.visibility)
	}
	if len(rule.when.predicates) != 1 {
		t.Fatalf("expected single predicate, got %d", len(rule.when.predicates))
	}

	p := rule.when.predicates[0]
	if p.metric != "messagesSent" || p.op != opGt || p.value != 1000 {
		t.Fatalf("unexpected predicate %+v", p)
	}
}

func TestMakeStoryRuleShouldDecodePayloadAndPredicate(t *testing.T) {
	payload := []byte(`{"id":"stat-earned","type":"stat","value":"moneyEarned"}`)

	rule, err := makeStoryRule("private", "moneyEarned", "gt", 0, true, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rule.visibility != visibilityPrivate {
		t.Fatalf("unexpected visibility %q", rule.visibility)
	}
	if rule.scene["id"] != "stat-earned" || rule.scene["value"] != "moneyEarned" {
		t.Fatalf("unexpected scene %+v", rule.scene)
	}
	if len(rule.when.predicates) != 1 || rule.when.predicates[0].metric != "moneyEarned" {
		t.Fatalf("unexpected predicate %+v", rule.when.predicates)
	}
}

func TestMakeStoryRuleWithoutPredicateShouldAlwaysMatch(t *testing.T) {
	rule, err := makeStoryRule("both", "", "", 0, false, []byte(`{"id":"intro-like"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rule.when.predicates) != 0 {
		t.Fatalf("expected no predicates, got %d", len(rule.when.predicates))
	}
	if !rule.when.eval(map[models.MetricKey]float64{}) {
		t.Fatal("predicate-less condition should always pass")
	}
}

func TestMakeStoryRuleShouldFailOnInvalidPayload(t *testing.T) {
	if _, err := makeStoryRule("both", "", "", 0, false, []byte(`{`)); err == nil {
		t.Fatal("expected error on malformed payload")
	}
}

func TestRuleProviderShouldFailWithoutStore(t *testing.T) {
	provider := NewRuleProvider(nil, time.Second)

	if _, err := provider.Get(context.Background()); err == nil {
		t.Fatal("expected error when rule store is missing")
	}
}
