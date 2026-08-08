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
- `GET /users/{userID}/metrics?from=&to=` — both bounds required together (RFC3339); if omitted, uses the current year in the user timezone from ClickHouse `users`. Each metric is an object `{ value, percentile, share }`. Fields are always present; `null` means no data or not applicable.

### Metric object semantics

| Category | value | share | percentile |
|---|---|---|---|
| counter | user total, `0` if no events | `user / global * 100` | % of users with a lower total |
| unique | user distinct count, `0` if empty | `user / global * 100` | % of users with a lower count |
| gauge | latest value, all `null` if absent | `null` | % of users with a lower gauge |
| milestone | first event unix time, all `null` if absent | `null` | % of users with a later first event |
| interval | avg duration seconds, all `null` if absent | `null` | % of users with a slower avg |

Example:

```json
{
  "metrics": {
    "listingsPublished": { "value": 47, "percentile": 88, "share": 4.7 },
    "firstListingAt": { "value": 1700000000, "percentile": 72, "share": null },
    "activeListings": { "value": null, "percentile": null, "share": null }
  }
}
```

## Local

Env vars: `CLICKHOUSE_HOST`, `CLICKHOUSE_PORT` (native, default `9000`), `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`, `CLICKHOUSE_DATABASE`, `SERVER_PORT`, `MIGRATIONS_DIR`.

Migrations run on startup from `migrations/`.
