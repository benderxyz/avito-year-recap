package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	conn driver.Conn
}

func Connect(ctx context.Context, host string, port int, user, password, database string) (*Client, error) {
	var lastErr error

	for attempt := 1; attempt <= 30; attempt++ {
		conn, err := clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", host, port)},
			Auth: clickhouse.Auth{
				Database: database,
				Username: user,
				Password: password,
			},
			DialTimeout: 5 * time.Second,
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
		})
		if err != nil {
			lastErr = fmt.Errorf("open clickhouse: %w", err)
		} else if err := conn.Ping(ctx); err != nil {
			_ = conn.Close()
			lastErr = fmt.Errorf("ping clickhouse: %w", err)
		} else {
			return &Client{conn: conn}, nil
		}

		fmt.Printf("clickhouse connect attempt %d/30 failed: %v\n", attempt, lastErr)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}

	return nil, fmt.Errorf("clickhouse unavailable after retries: %w", lastErr)
}

func (c *Client) Conn() driver.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Migrate(ctx context.Context, migrationsDir string) error {
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
			if err := c.conn.Exec(ctx, statement); err != nil {
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
