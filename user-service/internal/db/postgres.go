package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

const (
	connectAttempts   = 30
	connectRetryDelay = time.Second

	maxOpenConns    = 25
	maxIdleConns    = 5
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = time.Minute
)

func Connect(ctx context.Context, dsn string) (*Postgres, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetConnMaxLifetime(connMaxLifetime)
	conn.SetConnMaxIdleTime(connMaxIdleTime)

	var lastErr error

	for attempt := 1; attempt <= connectAttempts; attempt++ {
		lastErr = conn.PingContext(ctx)
		if lastErr == nil {
			return &Postgres{db: conn}, nil
		}

		slog.Warn(
			"postgres not ready, retrying",
			"attempt", attempt,
			"attempts", connectAttempts,
			"error", lastErr,
		)

		select {
		case <-ctx.Done():
			_ = conn.Close()

			return nil, ctx.Err()
		case <-time.After(connectRetryDelay):
		}
	}

	_ = conn.Close()

	return nil, fmt.Errorf("postgres unavailable after %d attempts: %w", connectAttempts, lastErr)
}

func (p *Postgres) DB() *sql.DB {
	return p.db
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

const createSchemaMigrations = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version    text PRIMARY KEY,
		applied_at timestamptz NOT NULL DEFAULT now()
	)
`

func (p *Postgres) Migrate(ctx context.Context, migrationsDir string) error {
	if _, err := p.db.ExecContext(ctx, createSchemaMigrations); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	applied, err := p.appliedMigrations(ctx)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		if _, ok := applied[name]; ok {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if err := p.applyMigration(ctx, name, string(content)); err != nil {
			return err
		}
	}

	return nil
}

func (p *Postgres) appliedMigrations(ctx context.Context) (map[string]struct{}, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	defer func() { _ = rows.Close() }()

	applied := make(map[string]struct{})
	for rows.Next() {
		var version string
		if scanErr := rows.Scan(&version); scanErr != nil {
			return nil, fmt.Errorf("read applied migrations: %w", scanErr)
		}
		applied[version] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}

	return applied, nil
}

func (p *Postgres) applyMigration(ctx context.Context, name string, content string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", name, err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, content); err != nil {
		return fmt.Errorf("exec migration %s: %w", name, err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO schema_migrations (version) VALUES ($1)`,
		name,
	); err != nil {
		return fmt.Errorf("record migration %s: %w", name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", name, err)
	}

	return nil
}
