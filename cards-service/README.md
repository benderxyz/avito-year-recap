# Cards Service

Builds the personalized "year recap" payload for a user. It fetches the user
profile from **user-service** and yearly metrics from **analytics-service**, runs
a data-driven rules engine (badges, story scenes, recommendations), and returns a
single `RecapPayload` consumed by the frontend recap widget.

An optional **LLM enrichment** step (off by default) can rewrite badge copy and
add one personality "insight" scene. It is strictly additive: on any failure the
base recap is served unchanged. See [LLM enrichment](#llm-enrichment).

## Key files

- [entrypoint](cmd/main.go)
- [HTTP handlers / routes](internal/api/handlers.go)
- [service clients (user / analytics)](internal/clients/clients.go), [LLM client](internal/clients/llm_client.go)
- [recap generator](internal/cards/generator.go)
- [rules engine](internal/cards/rules.go), [Postgres rule store](internal/cards/rulestore.go), [cached provider](internal/cards/rules_provider.go)
- [LLM enrichment service](internal/llm/service.go), [provider](internal/llm/openai.go), [prompt](internal/llm/prompt.go), [output validation](internal/llm/validate.go), [result cache](internal/llm/cache.go)
- [recap models](internal/models/recap.go)
- [admin routes](internal/admin/handler.go), [admin store](internal/admin/store_postgres.go), [validation](internal/admin/validate.go), [bearer auth](internal/admin/auth.go), [preview data](internal/admin/preview.go), [OpenAPI document](internal/admin/openapi.go)
- [schema DDL](migrations/001_rules.sql), [seed data](seeds/001_metric_definitions.sql)

## API

### `GET /health`

Health check. Response: `text/plain`, body `cards-service: OK`.

### `GET /api/recap/{year}/{id}`

Builds a **private** recap for the given user `id` and `year`.

- Fetches the profile from user-service and metrics from analytics-service, then
  runs the generator in `private` mode: full metrics, recommendations, and a
  share action (when a signing key is configured).
- When LLM enrichment is enabled, badge copy is rewritten and an insight scene is
  added on top of this payload (see [LLM enrichment](#llm-enrichment)).
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
- LLM enrichment, when enabled, applies here too — but only **public** metrics
  ever reach the provider (see [LLM enrichment](#llm-enrichment)).
- Success: `200 OK`, `application/json`, body `RecapPayload`.
- Errors:
  - `400 invalid share token` — token missing, malformed, or signature invalid.
  - `404 recap not found` — profile or metrics could not be resolved.

### `/api/admin/**`

Admin API over the four rule catalogs: metric definitions, badges, story scenes,
recommendations. See [Admin API](#admin-api).

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
    "moneyEarned": { "type": "money", "value": 150000, "currency": "RUB" },
    "firstListingAt": { "type": "date", "value": "2026-01-14" }
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
`outro`, and `insight` when LLM enrichment adds one). `features` is present only
in private mode when sharing is enabled.

## Rules & seeds

Badges, story scenes, recommendations, metric definitions come from Postgres.
Built-in fallback rules are gone. Postgres is required.

No metric key is hardcoded in Go. Every metric is a row in `metric_definitions`:

| column | meaning |
| --- | --- |
| `key` | key in the payload `metrics` map |
| `value_type` | `number`, `money`, `percentile`, `ratio`, `string`, `date` |
| `source_key` | key in the analytics-service response, defaults to `key` |
| `source_field` | which field of the analytics metric to read: `value`, `percentile`, `share` |
| `percentile_key` | metric that holds the percentile for comparison scenes |
| `comparison_min_percentile` | minimum percentile at which the comparison line is shown for this metric; empty means the service default of 50; percentile 0 is never shown regardless of the threshold |
| `is_public` | whether the metric survives in a public share card |
| `include_in_llm` | whether the metric goes into the LLM prompt |

A `date` metric is stored upstream as a unix timestamp and serialized as an ISO
date (`2026-01-14`). Formatting for humans happens in the frontend through
`Intl.DateTimeFormat`, so the `{{value}}` placeholder in scene copy is left
untouched by the backend and resolved by recap-engine using the scene `value`
key. Adding a metric means adding a row here plus a row in the analytics-service
`event_registry`, with no code change.

On startup the service migrates schema from `migrations/`, then loads rules
through a TTL-cached provider. Missing `POSTGRES_HOST`, connect failure, or
migrate failure stops the process.

Seed rows live in `seeds/` (separate from schema migrations). They apply on
startup only when `SEED_ON_START=true`, and only to empty tables. One file per
table, named `NNN_<table>.sql`: before running a file the service counts rows in
that table and skips the file when the table already has any. An empty stand
still gets the full catalog from `seeds/`; after the first start the admin API is
the source of truth, and rules edited there survive restarts.

## Admin API

`/api/admin/**` edits the same four tables the generator reads. Every request
needs `Authorization: Bearer <ADMIN_API_TOKEN>`; the token is compared in
constant time and never logged. An empty `ADMIN_API_TOKEN` is treated as a
configuration error: the service keeps serving the public recap and answers
`401` to every admin request.

Each catalog exposes the same five operations. `POST` creates (`409` when the id
is taken), `PUT` updates an existing record (`404` when it is missing) — they are
separate operations, not one upsert. Ids cannot be renamed: delete and create
instead.

| Method | Path | Meaning |
| --- | --- | --- |
| `GET` | `/api/admin/{catalog}` | list without pagination, including `enabled: false` |
| `POST` | `/api/admin/{catalog}` | create, `201` with the stored record |
| `GET` | `/api/admin/{catalog}/{id}` | one record |
| `PUT` | `/api/admin/{catalog}/{id}` | update, `200` with the stored record |
| `DELETE` | `/api/admin/{catalog}/{id}` | delete, `204` |

Catalogs are `metrics` (keyed by `key`), `badges`, `stories`, `recommendations`.
Lists are sorted the way the generator sorts them: `sortOrder` then id, and
`priority` descending for recommendations.

Filters are query parameters and are validated against the same enums as the
bodies, so an unknown value answers `400` instead of silently returning
everything:

| Catalog | Filters |
| --- | --- |
| `metrics` | `enabled`, `isPublic`, `includeInLlm`, `valueType`, `sourceField`, `search` |
| `badges` | `enabled`, `visibility`, `metric`, `search` |
| `stories` | `enabled`, `visibility`, `sceneType`, `metric`, `search` |
| `recommendations` | `enabled`, `metric`, `minPriority`, `search` |

Deleting a record that others point at answers `409` with the referencing ids in
`references`: a metric used by a badge, a scene, a recommendation predicate, or
another metric's `percentileKey` / `sourceKey`; a badge shown by a scene.

Errors always look like `{"error": "...", "fields": {...}, "references": [...]}`,
where `fields` maps a request field to the reason it was rejected.

`GET /api/admin/preview` returns a `RecapPayload` built by the same generator as
the public route, on throwaway random metrics derived from the current metric
definitions, and never calls the LLM. Every call differs; `?seed=` makes it
reproducible, `?year=` and `?mode=private|public` pick the recap being built.

`GET /api/admin/openapi.json` and `GET /api/admin/docs` (Swagger UI) describe all
of the above. Both are static API documentation and need no token, so the docs
page opens in a browser; requests fired from it still need the token, entered
through Swagger UI's Authorize button. Schemas are generated from the Go request
and response structs, so they cannot drift from the handlers.

Any write invalidates the rule provider cache, so the next recap and the next
preview see the change without waiting for the TTL.

## LLM enrichment

Optional post-processing that runs **after** the base recap is built. It rewrites
the selected badges' copy and may prepend one `insight` scene right after the
`intro`. It is enabled only when `LLM_ENABLED=true` and `OPENAI_API_KEY` is set;
otherwise the base recap is served unchanged.

**Strictly additive / fail-open.** Enrichment never blocks or breaks the recap.
On timeout, transport error, provider error, or unparsable output, the original
payload is returned and the failure is logged. The base recap flow is untouched.

**Privacy: only public metrics leave the service.** The prompt is built from
metric definitions that are both `is_public` and `include_in_llm`
([`RuleSet.MetricDefinitions()`](internal/cards/rules.go)) — the single source of
truth shared with share cards. Private metrics (e.g. `moneyEarned`, `moneySaved`,
`messagesSent`, `activeListings`) are **never** sent to the
provider, in either mode. The filter is **fail-closed**: if the allowlist is
empty, no metrics are sent at all. Model output is additionally validated
([validate.go](internal/llm/validate.go)) to strip digits, newlines, and
over-long text, so concrete numbers never appear in generated copy.

**Result cache.** Generations are cached in Postgres (`recap_llm_cache`) keyed on
`(user_id, year, mode, prompt_version)`:
- `mode` keeps public and private results separate — a private-mode generation is
  never served to a public share request.
- `prompt_version` is `sha256(system prompt + model)[:12]`. Editing the system
  prompt or switching `LLM_MODEL` changes the key, so stale copy is regenerated
  automatically with no manual cache flush.

The cache table is created by [migration 004](migrations/004_llm_cache_key.sql)
and shares the single Postgres connection with the rules store.

**Provider.** The client speaks the OpenAI-compatible `/responses` API, so
OpenAI and OpenRouter both work by pointing `LLM_BASE_URL` / `LLM_MODEL` at the
desired endpoint. The canonical setup is OpenRouter; `OPENAI_API_KEY` holds
whichever provider's key matches `LLM_BASE_URL`. Requests retry on `429`/`5xx`
(up to 3 attempts, short backoff) while respecting the request timeout.

## Configuration

Environment variables (see [internal/config/config.go](internal/config/config.go)):

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8081` | HTTP listen port |
| `LOG_LEVEL` | `info` | slog level: `debug`, `info`, `warn`, or `error` |
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
| `SEED_ON_START` | `false` | apply seeds from `SEEDS_DIR` on startup, skipping tables that already have rows |
| `ADMIN_API_TOKEN` | _(empty)_ | bearer token for `/api/admin/**`; empty closes the admin API with `401` |
| `LLM_ENABLED` | `false` | enable LLM enrichment (needs `OPENAI_API_KEY`) |
| `OPENAI_API_KEY` | _(empty)_ | provider API key matching `LLM_BASE_URL` (OpenRouter `sk-or-v1-…` or OpenAI `sk-…`) |
| `LLM_PROVIDER` | `openai` | provider adapter; `openai` covers any OpenAI-compatible endpoint |
| `LLM_BASE_URL` | `https://openrouter.ai/api/v1` | provider base URL (`/responses` is appended) |
| `LLM_MODEL` | `google/gemma-2-9b-it:free` | model slug (verify against your provider's catalog) |
| `LLM_TIMEOUT_MS` | `30000` | per-request timeout for the provider call |

Empty `POSTGRES_HOST`, connect failure, or migrate failure is fatal. The service
exits instead of serving recap without rules.

## CORS

The API is browser-facing, so [internal/api/cors.go](internal/api/cors.go) wraps
the mux with an origin allowlist from `CORS_ALLOWED_ORIGINS`. Origins are matched
exactly (scheme + host + port), so `http://localhost:3000` does not cover
`http://127.0.0.1:3000` — list both if you need both.

Allowed methods are `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS` and allowed headers
are `Content-Type` and `Authorization`, because the admin UI writes from the
browser. Writes stay protected by the admin token, not by the origin list.

Requests from a listed origin get `Access-Control-Allow-Origin`; preflight
`OPTIONS` is answered with `204` and never reaches the routes. Requests from an
unlisted origin are still served, just without CORS headers, so the browser
blocks them. Responses always carry `Vary: Origin` so a shared cache cannot serve
one origin's headers to another.

Credentials are not enabled: no `Access-Control-Allow-Credentials`, so browsers
will not send cookies cross-origin. Add it here (and drop any wildcard) if the
API ever moves to cookie auth.
