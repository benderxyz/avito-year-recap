package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, user User) error {
	timezone := user.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (user_id, external_id, username, timezone)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			external_id = EXCLUDED.external_id,
			username = EXCLUDED.username,
			timezone = EXCLUDED.timezone,
			updated_at = now()
	`, user.UserID, user.ExternalID, user.Username, timezone)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, userID uint64) (User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT user_id, external_id, username, timezone, created_at, updated_at
		FROM users
		WHERE user_id = $1
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

	return user, nil
}

func MapScanError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("get user: %w", err)
}
