package db

import (
	"context"
	"database/sql"
	"fmt"
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

func Connect(ctx context.Context, dsn string) (*Postgres, error) {
	var lastErr error

	for attempt := 1; attempt <= 30; attempt++ {
		conn, err := sql.Open("postgres", dsn)
		if err != nil {
			lastErr = fmt.Errorf("open postgres: %w", err)
		} else if err := conn.PingContext(ctx); err != nil {
			_ = conn.Close()
			lastErr = fmt.Errorf("ping postgres: %w", err)
		} else {
			return &Postgres{db: conn}, nil
		}

		fmt.Printf("postgres connect attempt %d/30 failed: %v\n", attempt, lastErr)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}

	return nil, fmt.Errorf("postgres unavailable after retries: %w", lastErr)
}

func (p *Postgres) DB() *sql.DB {
	return p.db
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) Migrate(ctx context.Context, migrationsDir string) error {
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
		path := filepath.Join(migrationsDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		statements := splitSQLStatements(string(content))
		for _, statement := range statements {
			if _, err := p.db.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("exec migration %s: %w", name, err)
			}
		}
	}

	return nil
}

func splitSQLStatements(content string) []string {
	parts := strings.Split(content, ";")
	statements := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		statements = append(statements, trimmed)
	}
	return statements
}
