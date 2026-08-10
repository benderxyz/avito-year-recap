# User Service

Manages test user profiles in Postgres. No event ingest — events go to analytics-service.

## API

- `GET /health`
- `PUT /users/{userID}` — upsert profile (`external_id`, `username`, `timezone`)
- `GET /users` — list all profiles
- `GET /users/{userID}` — read current profile

## Local

Env vars: `POSTGRES_HOST` (default `localhost`), `POSTGRES_PORT` (default `5432`), `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DATABASE` (default `users`), `POSTGRES_SSLMODE` (default `disable`), `SERVER_PORT`, `MIGRATIONS_DIR`.

In docker compose this service shares one Postgres container with cards-service and uses database `users`.

Migrations run on startup from `migrations/`. Upserts use `INSERT ... ON CONFLICT (user_id) DO UPDATE`; `created_at` / `updated_at` are managed by the database.
