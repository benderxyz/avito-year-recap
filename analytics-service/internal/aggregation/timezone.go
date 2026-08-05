package aggregation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseTimezoneResolver struct {
	conn driver.Conn
}

func NewClickHouseTimezoneResolver(conn driver.Conn) *ClickHouseTimezoneResolver {
	return &ClickHouseTimezoneResolver{conn: conn}
}

func (r *ClickHouseTimezoneResolver) Timezone(ctx context.Context, userID uint64) (string, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT timezone
		FROM users FINAL
		WHERE user_id = ?
		LIMIT 1
	`, userID)

	var timezone string
	if err := row.Scan(&timezone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "UTC", nil
		}
		return "", fmt.Errorf("resolve timezone: %w", err)
	}

	timezone = strings.TrimSpace(timezone)
	if timezone == "" {
		return "UTC", nil
	}

	return timezone, nil
}
