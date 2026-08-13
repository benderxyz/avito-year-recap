package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestSeedTableNameShouldExtractTableFromNumberedFile(t *testing.T) {
	got := seedTableName("003_story_rules.sql")

	if got != "story_rules" {
		t.Fatalf("want story_rules, got %q", got)
	}
}

func TestSeedTableNameShouldIgnoreFileWithoutNumberPrefix(t *testing.T) {
	got := seedTableName("story_rules.sql")

	if got != "" {
		t.Fatalf("want no table for unnumbered file, got %q", got)
	}
}

func TestSeedTableNameShouldIgnoreUppercaseTable(t *testing.T) {
	got := seedTableName("001_StoryRules.sql")

	if got != "" {
		t.Fatalf("want no table for unexpected file name, got %q", got)
	}
}

func newSeedTestPostgres(t *testing.T) (*Postgres, context.Context) {
	t.Helper()

	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN is not set; skipping database-backed tests")
	}

	ctx := context.Background()

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	if err := conn.PingContext(ctx); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}

	t.Cleanup(func() {
		if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS seed_probe"); err != nil {
			t.Fatalf("drop probe table: %v", err)
		}
		_ = conn.Close()
	})

	if _, err := conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS seed_probe (id text PRIMARY KEY)"); err != nil {
		t.Fatalf("create probe table: %v", err)
	}

	return &Postgres{db: conn}, ctx
}

func writeSeedFile(t *testing.T, statement string) string {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "001_seed_probe.sql"), []byte(statement), 0o600); err != nil {
		t.Fatalf("write seed file: %v", err)
	}

	return dir
}

func TestSeedShouldFillEmptyTable(t *testing.T) {
	pg, ctx := newSeedTestPostgres(t)
	dir := writeSeedFile(t, "INSERT INTO seed_probe (id) VALUES ('from-seed');")

	if err := pg.Seed(ctx, dir); err != nil {
		t.Fatalf("seed empty table: %v", err)
	}

	var stored string
	if err := pg.db.QueryRowContext(ctx, "SELECT id FROM seed_probe").Scan(&stored); err != nil {
		t.Fatalf("read probe table: %v", err)
	}
	if stored != "from-seed" {
		t.Fatalf("want row from seed file, got %q", stored)
	}
}

func TestSeedShouldKeepRowsEditedInAdmin(t *testing.T) {
	pg, ctx := newSeedTestPostgres(t)
	if _, err := pg.db.ExecContext(ctx, "INSERT INTO seed_probe (id) VALUES ('from-admin')"); err != nil {
		t.Fatalf("insert admin row: %v", err)
	}
	dir := writeSeedFile(t, "INSERT INTO seed_probe (id) VALUES ('from-seed');")

	if err := pg.Seed(ctx, dir); err != nil {
		t.Fatalf("seed non empty table: %v", err)
	}

	var rows int
	if err := pg.db.QueryRowContext(ctx, "SELECT count(*) FROM seed_probe").Scan(&rows); err != nil {
		t.Fatalf("count probe rows: %v", err)
	}
	if rows != 1 {
		t.Fatalf("want seed skipped for non empty table, got %d rows", rows)
	}
}
