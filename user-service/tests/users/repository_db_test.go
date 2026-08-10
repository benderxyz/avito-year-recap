package users_test

import (
	"context"
	"errors"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"user-service/internal/db"
	"user-service/internal/users"
)

func newTestRepository(t *testing.T) (*users.Repository, context.Context) {
	t.Helper()

	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN is not set; skipping database-backed tests")
	}

	ctx := context.Background()

	pg, err := db.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("connect test postgres: %v", err)
	}
	t.Cleanup(func() { _ = pg.Close() })

	if err := pg.Migrate(ctx, "../../migrations"); err != nil {
		t.Fatalf("migrate test postgres: %v", err)
	}

	if _, err := pg.DB().ExecContext(ctx, `TRUNCATE users`); err != nil {
		t.Fatalf("truncate users: %v", err)
	}

	return users.NewRepository(pg.DB()), ctx
}

func TestUpsertShouldReturnStoredRow(t *testing.T) {
	repo, ctx := newTestRepository(t)

	saved, err := repo.Upsert(ctx, users.User{
		UserID:     42,
		ExternalID: "avito-42",
		Username:   "alex",
		Timezone:   "Europe/Moscow",
	})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}

	if saved.UserID != 42 || saved.Username != "alex" || saved.Timezone != "Europe/Moscow" {
		t.Fatalf("unexpected stored row: %+v", saved)
	}
	if saved.CreatedAt.IsZero() || saved.UpdatedAt.IsZero() {
		t.Fatalf("expected timestamps to be populated, got %+v", saved)
	}
}

func TestUpsertShouldDefaultEmptyTimezone(t *testing.T) {
	repo, ctx := newTestRepository(t)

	saved, err := repo.Upsert(ctx, users.User{UserID: 7})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}

	if saved.Timezone != users.DefaultTimezone {
		t.Fatalf("expected %q, got %q", users.DefaultTimezone, saved.Timezone)
	}
}

func TestUpsertShouldOverwriteExistingUser(t *testing.T) {
	repo, ctx := newTestRepository(t)

	first, err := repo.Upsert(ctx, users.User{UserID: 42, Username: "before", Timezone: "UTC"})
	if err != nil {
		t.Fatalf("first upsert: %v", err)
	}

	second, err := repo.Upsert(ctx, users.User{UserID: 42, Username: "after", Timezone: "Europe/Moscow"})
	if err != nil {
		t.Fatalf("second upsert: %v", err)
	}

	if second.Username != "after" || second.Timezone != "Europe/Moscow" {
		t.Fatalf("expected updated values, got %+v", second)
	}
	if !second.CreatedAt.Equal(first.CreatedAt) {
		t.Fatal("expected created_at to be preserved on update")
	}
	if second.UpdatedAt.Before(first.UpdatedAt) {
		t.Fatal("expected updated_at to move forward")
	}
}

func TestUpsertShouldRejectDuplicateExternalID(t *testing.T) {
	repo, ctx := newTestRepository(t)

	if _, err := repo.Upsert(ctx, users.User{UserID: 1, ExternalID: "avito-1"}); err != nil {
		t.Fatalf("first upsert: %v", err)
	}

	_, err := repo.Upsert(ctx, users.User{UserID: 2, ExternalID: "avito-1"})
	if !errors.Is(err, users.ErrExternalIDTaken) {
		t.Fatalf("expected ErrExternalIDTaken, got %v", err)
	}
}

func TestUpsertShouldAllowMultipleEmptyExternalIDs(t *testing.T) {
	repo, ctx := newTestRepository(t)

	if _, err := repo.Upsert(ctx, users.User{UserID: 1}); err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	if _, err := repo.Upsert(ctx, users.User{UserID: 2}); err != nil {
		t.Fatalf("expected empty external_id to be exempt from the unique index, got %v", err)
	}
}

func TestGetByIDShouldReturnNotFoundForMissingUser(t *testing.T) {
	repo, ctx := newTestRepository(t)

	if _, err := repo.GetByID(ctx, 404); !errors.Is(err, users.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListShouldReturnUsersOrderedByID(t *testing.T) {
	repo, ctx := newTestRepository(t)

	for _, id := range []uint64{3, 1, 2} {
		if _, err := repo.Upsert(ctx, users.User{UserID: id}); err != nil {
			t.Fatalf("seed user %d: %v", id, err)
		}
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(list) != 3 {
		t.Fatalf("expected 3 users, got %d", len(list))
	}
	for i, want := range []uint64{1, 2, 3} {
		if list[i].UserID != want {
			t.Fatalf("expected user %d at index %d, got %d", want, i, list[i].UserID)
		}
	}
}

func TestListShouldReturnEmptySliceWhenNoUsers(t *testing.T) {
	repo, ctx := newTestRepository(t)

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if list == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(list) != 0 {
		t.Fatalf("expected no users, got %d", len(list))
	}
}
