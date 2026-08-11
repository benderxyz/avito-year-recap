package events_test

import (
	"testing"

	"analytics-service/internal/events"
)

func TestRegistryShouldReturnConfigWhenEventTypeIsRegistered(t *testing.T) {
	registry := testRegistry()

	cfg, err := registry.Get("item_published")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Category != events.CategoryCounter {
		t.Fatalf("expected category %s, got %s", events.CategoryCounter, cfg.Category)
	}
	if cfg.MetricKey != "listingsPublished" {
		t.Fatalf("expected metric key listingsPublished, got %s", cfg.MetricKey)
	}
}

func TestRegistryShouldFailWhenEventTypeIsUnknown(t *testing.T) {
	registry := testRegistry()

	_, err := registry.Get("not_registered")

	if err == nil {
		t.Fatal("expected error for unknown event type")
	}
}

func TestRegistryShouldExposeUniqueFieldForCategoryOpened(t *testing.T) {
	registry := testRegistry()

	cfg, err := registry.Get("category_opened")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Category != events.CategoryUnique {
		t.Fatalf("expected category %s, got %s", events.CategoryUnique, cfg.Category)
	}
	if cfg.UniqueField != "category" {
		t.Fatalf("expected unique field category, got %s", cfg.UniqueField)
	}
}

func TestRegistryShouldUseDayModeForDayActive(t *testing.T) {
	registry := testRegistry()

	cfg, err := registry.Get("day_active")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.UniqueMode != events.UniqueModeDay {
		t.Fatalf("expected unique mode day, got %s", cfg.UniqueMode)
	}
	if cfg.MetricKey != "daysActive" {
		t.Fatalf("expected metric key daysActive, got %s", cfg.MetricKey)
	}
}
