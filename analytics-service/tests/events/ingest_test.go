package events_test

import (
	"testing"
	"time"

	"analytics-service/internal/events"
)

func TestNormalizeShouldAssignCategoryFromRegistryWhenEventTypeIsValid(t *testing.T) {
	event, err := events.BuildEvent(testRegistry(), events.IngestEvent{
		UserID:    42,
		EventType: "item_view_received",
		Value:     10,
	}, time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.EventCategory != events.CategoryCounter {
		t.Fatalf("expected counter category, got %s", event.EventCategory)
	}
	if event.Value != 10 {
		t.Fatalf("expected counter value 10, got %v", event.Value)
	}
	if event.Payload != "{}" {
		t.Fatalf("expected default payload {}, got %s", event.Payload)
	}
}

func TestBuildEventShouldForceZeroValueWhenCategoryIsMilestone(t *testing.T) {
	event, err := events.BuildEvent(testRegistry(), events.IngestEvent{
		UserID:    42,
		EventType: "first_item_published",
		Value:     99,
	}, time.Now().UTC())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.EventCategory != events.CategoryMilestone {
		t.Fatalf("expected milestone category, got %s", event.EventCategory)
	}
	if event.Value != 0 {
		t.Fatalf("expected milestone value 0, got %v", event.Value)
	}
}

func TestBuildEventShouldForceZeroValueWhenCategoryIsUnique(t *testing.T) {
	event, err := events.BuildEvent(testRegistry(), events.IngestEvent{
		UserID:    42,
		EventType: "category_opened",
		Value:     5,
		Payload:   `{"category":"electronics"}`,
	}, time.Now().UTC())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Value != 0 {
		t.Fatalf("expected unique value 0, got %v", event.Value)
	}
}

func TestNormalizeShouldFailWhenUserIDIsMissing(t *testing.T) {
	_, err := events.BuildEvent(testRegistry(), events.IngestEvent{
		EventType: "item_view_received",
	}, time.Now().UTC())

	if err == nil {
		t.Fatal("expected error when user_id is missing")
	}
}
