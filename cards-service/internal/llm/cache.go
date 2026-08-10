package llm

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type Cache struct {
	db *sql.DB
}

func NewCache(db *sql.DB) *Cache {
	return &Cache{db: db}
}

func (c *Cache) Get(ctx context.Context, userID string, year int, mode, version string) (Result, bool, error) {
	if c == nil || c.db == nil {
		return Result{}, false, nil
	}

	var raw []byte
	err := c.db.QueryRowContext(ctx, `
		SELECT result FROM recap_llm_cache
		WHERE user_id = $1 AND year = $2 AND mode = $3 AND prompt_version = $4
	`, userID, year, mode, version).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return Result{}, false, nil
	}
	if err != nil {
		return Result{}, false, fmt.Errorf("query llm cache: %w", err)
	}

	var result Result
	if err := json.Unmarshal(raw, &result); err != nil {
		return Result{}, false, fmt.Errorf("unmarshal cached result: %w", err)
	}
	return result, true, nil
}

func (c *Cache) Put(ctx context.Context, userID string, year int, mode, version string, result Result) error {
	if c == nil || c.db == nil {
		return nil
	}

	raw, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	_, err = c.db.ExecContext(ctx, `
		INSERT INTO recap_llm_cache (user_id, year, mode, prompt_version, result)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, year, mode, prompt_version)
		DO UPDATE SET result = EXCLUDED.result, created_at = now()
	`, userID, year, mode, version, raw)
	if err != nil {
		return fmt.Errorf("upsert llm cache: %w", err)
	}
	return nil
}
