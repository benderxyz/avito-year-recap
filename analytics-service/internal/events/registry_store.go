package events

import (
	"context"
	"database/sql"
	"fmt"
)

type RegistryStore struct {
	db *sql.DB
}

func NewRegistryStore(db *sql.DB) *RegistryStore {
	return &RegistryStore{db: db}
}

func (s *RegistryStore) Load(ctx context.Context) (*Registry, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT event_type, category, metric_key, unique_mode, unique_field
		FROM event_registry
		WHERE enabled = true
		ORDER BY sort_order, event_type
	`)
	if err != nil {
		return nil, fmt.Errorf("query event_registry: %w", err)
	}
	defer func() { _ = rows.Close() }()

	configs := make(map[string]CategoryConfig)
	for rows.Next() {
		var eventType, category, metricKey string
		var uniqueMode, uniqueField sql.NullString
		if err := rows.Scan(&eventType, &category, &metricKey, &uniqueMode, &uniqueField); err != nil {
			return nil, fmt.Errorf("scan event_registry: %w", err)
		}
		configs[eventType] = CategoryConfig{
			Category:    EventCategory(category),
			MetricKey:   metricKey,
			UniqueMode:  UniqueMode(uniqueMode.String),
			UniqueField: uniqueField.String,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate event_registry: %w", err)
	}

	return NewRegistry(configs), nil
}
