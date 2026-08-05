package events

import (
	"time"

	"github.com/google/uuid"
)

type EventCategory string

const (
	CategoryCounter   EventCategory = "counter"
	CategoryInterval  EventCategory = "interval"
	CategoryGauge     EventCategory = "gauge"
	CategoryMilestone EventCategory = "milestone"
	CategoryUnique    EventCategory = "unique"
)

type Event struct {
	EventID       uuid.UUID     `json:"event_id"`
	UserID        uint64        `json:"user_id"`
	SessionID     uuid.UUID     `json:"session_id"`
	EventType     string        `json:"event_type"`
	EventCategory EventCategory `json:"event_category"`
	Value         float64       `json:"value"`
	Payload       string        `json:"payload"`
	OccurredAt    time.Time     `json:"occurred_at"`
	InsertedAt    time.Time     `json:"inserted_at"`
}

type IngestEvent struct {
	UserID     uint64    `json:"user_id"`
	SessionID  string    `json:"session_id"`
	EventType  string    `json:"event_type"`
	Value      float64   `json:"value"`
	Payload    string    `json:"payload"`
	OccurredAt time.Time `json:"occurred_at"`
}

func ValueForCategory(category EventCategory, value float64) float64 {
	switch category {
	case CategoryMilestone, CategoryUnique:
		return 0
	default:
		return value
	}
}
