package users_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"user-service/internal/users"
)

func TestApplyUpsertTimestampsShouldPreserveCreatedAtWhenUserExists(t *testing.T) {
	createdAt := time.Date(2025, 1, 10, 8, 0, 0, 0, time.UTC)
	now := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	existing := users.User{UserID: 42, CreatedAt: createdAt}

	got := users.ApplyUpsertTimestamps(users.User{
		UserID:   42,
		Username: "Alex",
		Timezone: "Europe/Moscow",
	}, &existing, now)

	if !got.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected created_at %v, got %v", createdAt, got.CreatedAt)
	}
	if !got.UpdatedAt.Equal(now) {
		t.Fatalf("expected updated_at %v, got %v", now, got.UpdatedAt)
	}
}

func TestApplyUpsertTimestampsShouldSetCreatedAtWhenUserIsNew(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)

	got := users.ApplyUpsertTimestamps(users.User{UserID: 7}, nil, now)

	if !got.CreatedAt.Equal(now) {
		t.Fatalf("expected created_at %v, got %v", now, got.CreatedAt)
	}
	if got.Timezone != "UTC" {
		t.Fatalf("expected default timezone UTC, got %s", got.Timezone)
	}
}

func TestMapScanErrorShouldReturnNotFoundWhenNoRows(t *testing.T) {
	err := users.MapScanError(sql.ErrNoRows)

	if !errors.Is(err, users.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMapScanErrorShouldWrapUnexpectedErrors(t *testing.T) {
	original := errors.New("connection reset")

	err := users.MapScanError(original)

	if errors.Is(err, users.ErrNotFound) {
		t.Fatal("did not expect ErrNotFound for connection error")
	}
	if !errors.Is(err, original) {
		t.Fatalf("expected wrapped original error, got %v", err)
	}
	if fmt.Sprint(err) == original.Error() {
		t.Fatal("expected wrapped error message")
	}
}
