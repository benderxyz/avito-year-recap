package events

import (
	"context"
	"fmt"
	"strings"
	"time"

	"analytics-service/internal/apperr"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
)

type Ingester struct {
	conn     driver.Conn
	registry *Registry
}

func NewIngester(conn driver.Conn, registry *Registry) *Ingester {
	return &Ingester{conn: conn, registry: registry}
}

func (i *Ingester) Ingest(ctx context.Context, items []IngestEvent) error {
	if len(items) == 0 {
		return apperr.Validation("empty events batch")
	}

	events := make([]Event, 0, len(items))
	now := time.Now().UTC()

	for idx, item := range items {
		event, err := i.BuildEvent(item, now)
		if err != nil {
			return apperr.Validationf("event[%d]: %v", idx, err)
		}
		events = append(events, event)
	}

	batch, err := i.conn.PrepareBatch(ctx, `
		INSERT INTO events (
			event_id, user_id, session_id, event_type, event_category,
			value, payload, occurred_at, inserted_at
		)
	`)
	if err != nil {
		return fmt.Errorf("prepare events batch: %w", err)
	}

	for _, event := range events {
		if err := batch.Append(
			event.EventID,
			event.UserID,
			event.SessionID,
			event.EventType,
			string(event.EventCategory),
			event.Value,
			event.Payload,
			event.OccurredAt,
			event.InsertedAt,
		); err != nil {
			return fmt.Errorf("append event: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("send events batch: %w", err)
	}

	return nil
}

func (i *Ingester) BuildEvent(item IngestEvent, now time.Time) (Event, error) {
	if item.UserID == 0 {
		return Event{}, fmt.Errorf("user_id is required")
	}
	if strings.TrimSpace(item.EventType) == "" {
		return Event{}, fmt.Errorf("event_type is required")
	}

	cfg, err := i.registry.Get(item.EventType)
	if err != nil {
		return Event{}, err
	}

	sessionID := uuid.Nil
	if strings.TrimSpace(item.SessionID) != "" {
		parsed, err := uuid.Parse(item.SessionID)
		if err != nil {
			return Event{}, fmt.Errorf("invalid session_id")
		}
		sessionID = parsed
	}

	payload := item.Payload
	if strings.TrimSpace(payload) == "" {
		payload = "{}"
	}

	occurredAt := item.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = now
	}

	return Event{
		EventID:       uuid.New(),
		UserID:        item.UserID,
		SessionID:     sessionID,
		EventType:     item.EventType,
		EventCategory: cfg.Category,
		Value:         ValueForCategory(cfg.Category, item.Value),
		Payload:       payload,
		OccurredAt:    occurredAt.UTC(),
		InsertedAt:    now,
	}, nil
}
