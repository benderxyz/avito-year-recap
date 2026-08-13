package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
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
	files, err := sqlFiles(migrationsDir)
	if err != nil {
		return err
	}

	for _, name := range files {
		if err := p.execSQLFile(ctx, migrationsDir, name); err != nil {
			return err
		}
	}

	return nil
}

var seedFileTable = regexp.MustCompile(`^\d+_([a-z0-9_]+)\.sql$`)

func (p *Postgres) Seed(ctx context.Context, seedsDir string) error {
	files, err := sqlFiles(seedsDir)
	if err != nil {
		return err
	}

	for _, name := range files {
		table := seedTableName(name)
		if table != "" {
			empty, err := p.tableIsEmpty(ctx, table)
			if err != nil {
				return err
			}
			if !empty {
				slog.Info("seed skipped, table already has rows", "file", name, "table", table)
				continue
			}
		}

		if err := p.execSQLFile(ctx, seedsDir, name); err != nil {
			return err
		}
	}

	return nil
}

func seedTableName(fileName string) string {
	match := seedFileTable.FindStringSubmatch(fileName)
	if match == nil {
		return ""
	}

	return match[1]
}

func (p *Postgres) tableIsEmpty(ctx context.Context, table string) (bool, error) {
	var relation sql.NullString
	if err := p.db.QueryRowContext(ctx, "SELECT to_regclass($1)", table).Scan(&relation); err != nil {
		return false, fmt.Errorf("check table %s: %w", table, err)
	}
	if !relation.Valid {
		return true, nil
	}

	var hasRows bool
	query := "SELECT EXISTS (SELECT 1 FROM " + pq.QuoteIdentifier(table) + ")"
	if err := p.db.QueryRowContext(ctx, query).Scan(&hasRows); err != nil {
		return false, fmt.Errorf("count rows in %s: %w", table, err)
	}

	return !hasRows, nil
}

func sqlFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read sql dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	return files, nil
}

func (p *Postgres) execSQLFile(ctx context.Context, dir, name string) error {
	content, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return fmt.Errorf("read sql file %s: %w", name, err)
	}

	for _, statement := range splitSQLStatements(string(content)) {
		if _, err := p.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("exec sql file %s: %w", name, err)
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
