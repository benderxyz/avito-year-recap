# Analytics Service

Stores user events in ClickHouse, aggregates them by category, and exposes metrics for recap.

## Categories

| Category | Meaning | Aggregation | `value` |
|---|---|---|---|
| `counter` | repeating numeric events | `sum(value)` | amount / count |
| `interval` | duration between start and end | `avg(end - start)` | unix start time |
| `gauge` | latest snapshot in the period | `argMax(value, occurred_at)` | current level |
| `milestone` | happened or not (0/1 time) | `Present` + `min(occurred_at)` | always stored as `0` |
| `unique` | distinct set members | `uniqExact(...)` | always stored as `0` |

For `unique`

- payload mode reads `JSONExtractString(payload, field)`
- day mode counts distinct local calendar days via user timezone

Event types are Avito marketplace actions (`item_published`, `deal_completed`, `category_opened`, …). See `internal/events/registry.go`.

## API

- `GET /health`
- `POST /events` — single event object or JSON array
- `GET /users/{userID}/metrics?from=&to=` — both bounds required together (RFC3339); if omitted, uses the current year in the user timezone from ClickHouse `users`. Sparse metrics (`gauge` / `interval` / `milestone`) are JSON `null` when there is no data.

## Local

Env vars: `CLICKHOUSE_HOST`, `CLICKHOUSE_PORT` (native, default `9000`), `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`, `CLICKHOUSE_DATABASE`, `SERVER_PORT`, `MIGRATIONS_DIR`.

Migrations run on startup from `migrations/`.
