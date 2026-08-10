package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrNotFound        = errors.New("user not found")
	ErrExternalIDTaken = errors.New("external_id already belongs to another user")
)

const uniqueViolation = "23505"

const userColumns = `user_id, external_id, username, timezone, created_at, updated_at`

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, user User) (User, error) {
	timezone := user.Timezone
	if timezone == "" {
		timezone = DefaultTimezone
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO users (user_id, external_id, username, timezone)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			external_id = EXCLUDED.external_id,
			username = EXCLUDED.username,
			timezone = EXCLUDED.timezone,
			updated_at = now()
		RETURNING `+userColumns,
		user.UserID, user.ExternalID, user.Username, timezone,
	)

	saved, err := scanUser(row)
	if err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrExternalIDTaken
		}

		return User{}, fmt.Errorf("upsert user: %w", err)
	}

	return saved, nil
}

func (r *Repository) GetByID(ctx context.Context, userID uint64) (User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+userColumns+`
		FROM users
		WHERE user_id = $1
	`, userID)

	user, err := scanUser(row)
	if err != nil {
		return User{}, MapScanError(err)
	}

	return user, nil
}

func (r *Repository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+userColumns+`
		FROM users
		ORDER BY user_id
	`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer func() { _ = rows.Close() }()

	users := make([]User, 0)
	for rows.Next() {
		user, scanErr := scanUser(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("list users: %w", scanErr)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row rowScanner) (User, error) {
	var user User
	err := row.Scan(
		&user.UserID,
		&user.ExternalID,
		&user.Username,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == uniqueViolation
	}

	return false
}

func MapScanError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	return fmt.Errorf("get user: %w", err)
}
