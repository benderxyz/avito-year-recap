# User Service

Manages test user profiles in ClickHouse (`ReplacingMergeTree`). No event ingest — events go to analytics-service.

## API

- `GET /health`
- `PUT /users/{userID}` — upsert profile (`external_id`, `username`, `timezone`)
- `GET /users/{userID}` — read current profile (`FINAL`)

## Local

Env vars: `CLICKHOUSE_HOST`, `CLICKHOUSE_PORT` (native, default `9000`), `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`, `CLICKHOUSE_DATABASE`, `SERVER_PORT`, `MIGRATIONS_DIR`.

Migrations run on startup from `migrations/`.
