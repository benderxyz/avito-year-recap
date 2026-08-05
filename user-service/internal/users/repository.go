package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	conn driver.Conn
}

func NewRepository(conn driver.Conn) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) Upsert(ctx context.Context, user User) error {
	now := time.Now().UTC()

	existing, err := r.GetByID(ctx, user.UserID)
	switch {
	case err == nil:
		user = ApplyUpsertTimestamps(user, &existing, now)
	case errors.Is(err, ErrNotFound):
		user = ApplyUpsertTimestamps(user, nil, now)
	default:
		return fmt.Errorf("lookup user before upsert: %w", err)
	}

	batch, err := r.conn.PrepareBatch(ctx, `
		INSERT INTO users (user_id, external_id, username, timezone, created_at, updated_at)
	`)
	if err != nil {
		return fmt.Errorf("prepare upsert batch: %w", err)
	}

	if err := batch.Append(
		user.UserID,
		user.ExternalID,
		user.Username,
		user.Timezone,
		user.CreatedAt,
		user.UpdatedAt,
	); err != nil {
		return fmt.Errorf("append user: %w", err)
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("send user upsert: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, userID uint64) (User, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT user_id, external_id, username, timezone, created_at, updated_at
		FROM users FINAL
		WHERE user_id = ?
		LIMIT 1
	`, userID)

	var user User
	if err := row.Scan(
		&user.UserID,
		&user.ExternalID,
		&user.Username,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return User{}, MapScanError(err)
	}

	if user.UserID == 0 {
		return User{}, ErrNotFound
	}

	return user, nil
}

func MapScanError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("get user: %w", err)
}
