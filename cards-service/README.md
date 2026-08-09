# Cards Service

Builds the personalized "year recap" payload for a user. It fetches the user
profile from **user-service** and yearly metrics from **analytics-service**, runs
a data-driven rules engine (badges, story scenes, recommendations), and returns a
single `RecapPayload` consumed by the frontend recap widget.

## Key files

- [entrypoint](cmd/main.go)
- [HTTP handlers / routes](internal/api/handlers.go)
- [service clients (user / analytics)](internal/clients/clients.go)
- [recap generator](internal/cards/generator.go)
- [rules engine](internal/cards/rules.go), [Postgres rule store](internal/cards/rulestore.go), [cached provider](internal/cards/rules_provider.go)
- [recap models](internal/models/recap.go)
- [schema DDL](migrations/001_rules.sql), [seed data](seeds/001_rules.sql)

## API

### `GET /health`

Health check. Response: `text/plain`, body `cards-service: OK`.

### `GET /api/recap/{year}/{id}`

Builds a **private** recap for the given user `id` and `year`.

- Fetches the profile from user-service and metrics from analytics-service, then
  runs the generator in `private` mode: full metrics, recommendations, and a
  share action (when a signing key is configured).
- Success: `200 OK`, `application/json`, body `RecapPayload` (see below).
- Errors:
  - `400 invalid year` — `year` is not an integer.
  - `404 recap not found` — profile or metrics could not be resolved.

### `GET /api/share/{token}`

Builds a **public** recap from a signed share token (as embedded in a private
recap's `features.shareUrl`).

- Decodes `id` and `year` from the token, then runs the generator in `public`
  mode: metrics limited to a public allowlist, no recommendations, no share
  action.
- Success: `200 OK`, `application/json`, body `RecapPayload`.
- Errors:
  - `400 invalid share token` — token missing, malformed, or signature invalid.
  - `404 recap not found` — profile or metrics could not be resolved.

## Response shape (`RecapPayload`)

```json
{
  "schemaVersion": 1,
  "meta": {
    "vertical": "marketplace",
    "year": 2026,
    "locale": "ru-RU",
    "user": { "id": "avito-42", "displayName": "Alex", "avatarUrl": "…" },
    "generatedAt": "2026-08-09T00:00:00Z"
  },
  "metrics": {
    "listingsPublished": { "type": "number", "value": 12 },
    "moneyEarned": { "type": "money", "value": 150000, "currency": "RUB" }
  },
  "badges": [
    { "id": "active_user", "title": "…", "description": "…", "iconUrl": "/badges/active_user.svg" }
  ],
  "story": [
    { "id": "intro", "type": "intro", "title": "…", "actions": [] }
  ],
  "features": { "shareEnabled": true, "shareUrl": "…" }
}
```

`story` is an ordered list of scenes (`intro`, `stat`, `achievement`, `upsell`,
`outro`). `features` is present only in private mode when sharing is enabled.

## Rules & seeds

Badges, story scenes, recommendations, metric definitions come from Postgres.
Built-in fallback rules are gone. Postgres is required.

On startup the service migrates schema from `migrations/`, then loads rules
through a TTL-cached provider. Missing `POSTGRES_HOST`, connect failure, or
migrate failure stops the process.

Seed rows live in `seeds/` (separate from schema migrations). They apply on
startup only when `SEED_ON_START=true`. Seeds use `ON CONFLICT DO NOTHING`, so
restarts stay idempotent.

## Configuration

Environment variables (see [internal/config/config.go](internal/config/config.go)):

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8081` | HTTP listen port |
| `USER_SERVICE_URL` | `http://localhost:8082` | user-service base URL |
| `ANALYTICS_SERVICE_URL` | `http://localhost:8080` | analytics-service base URL |
| `SHARE_SIGNING_KEY` | `dev-insecure-share-key` | HMAC key for share tokens |
| `SHARE_BASE_URL` | `http://localhost:3000` | base URL used to build share links |
| `PRODUCT_BASE_URL` | `https://www.avito.ru` | base URL for recommendation links |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3000` | comma-separated origins allowed to call the API from a browser |
| `POSTGRES_HOST` | _(empty)_ | Postgres host (required). Empty value stops the process. In compose shares one Postgres with user-service |
| `POSTGRES_PORT` | `5432` | Postgres port |
| `POSTGRES_USER` | `recap` | Postgres user |
| `POSTGRES_PASSWORD` | `recap` | Postgres password |
| `POSTGRES_DATABASE` | `cards` | Postgres database (`cards`, separate from `users`) |
| `POSTGRES_SSLMODE` | `disable` | Postgres SSL mode |
| `MIGRATIONS_DIR` | `migrations` | schema migrations directory |
| `SEEDS_DIR` | `seeds` | seed files directory |
| `SEED_ON_START` | `false` | apply seeds from `SEEDS_DIR` on startup |

Empty `POSTGRES_HOST`, connect failure, or migrate failure is fatal. The service
exits instead of serving recap without rules.

## CORS

The API is browser-facing, so [internal/api/cors.go](internal/api/cors.go) wraps
the mux with an origin allowlist from `CORS_ALLOWED_ORIGINS`. Origins are matched
exactly (scheme + host + port), so `http://localhost:3000` does not cover
`http://127.0.0.1:3000` — list both if you need both.

Requests from a listed origin get `Access-Control-Allow-Origin`; preflight
`OPTIONS` is answered with `204` and never reaches the routes. Requests from an
unlisted origin are still served, just without CORS headers, so the browser
blocks them. Responses always carry `Vary: Origin` so a shared cache cannot serve
one origin's headers to another.

Credentials are not enabled: no `Access-Control-Allow-Credentials`, so browsers
will not send cookies cross-origin. Add it here (and drop any wildcard) if the
API ever moves to cookie auth.
