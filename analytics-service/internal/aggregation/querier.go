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

func (q *ClickHouseQuerier) QueryFloat64s(ctx context.Context, query string, args ...any) ([]float64, error) {
	rows, err := q.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query float64s: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	totals := make([]float64, 0)
	for rows.Next() {
		var value sql.NullFloat64
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("scan float64: %w", err)
		}
		if !value.Valid || math.IsNaN(value.Float64) || math.IsInf(value.Float64, 0) {
			continue
		}
		totals = append(totals, sanitizeFloat(value.Float64))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return totals, nil
}

func sanitizeFloat(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return value
}
