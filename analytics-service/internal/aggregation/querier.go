package aggregation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseQuerier struct {
	conn driver.Conn
}

func NewClickHouseQuerier(conn driver.Conn) *ClickHouseQuerier {
	return &ClickHouseQuerier{conn: conn}
}

func (q *ClickHouseQuerier) QueryFloat64(ctx context.Context, query string, args ...any) (float64, bool, error) {
	row := q.conn.QueryRow(ctx, query, args...)

	var value sql.NullFloat64
	if err := row.Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("scan float64: %w", err)
	}

	if !value.Valid {
		return 0, false, nil
	}

	sanitized := sanitizeFloat(value.Float64)
	if math.IsNaN(value.Float64) || math.IsInf(value.Float64, 0) {
		return 0, false, nil
	}

	return sanitized, true, nil
}

func sanitizeFloat(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return value
}
