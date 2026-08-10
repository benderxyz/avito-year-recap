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
- [schema DDL](migrations/001_rules.sql), [seed data](seeds/001_rules.sql)

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
`outro`, and `insight` when LLM enrichment adds one). `features` is present only
in private mode when sharing is enabled.

## Rules & seeds

Badges, story scenes, recommendations, metric definitions come from Postgres.
Built-in fallback rules are gone. Postgres is required.

On startup the service migrates schema from `migrations/`, then loads rules
through a TTL-cached provider. Missing `POSTGRES_HOST`, connect failure, or
migrate failure stops the process.

Seed rows live in `seeds/` (separate from schema migrations). They apply on
startup only when `SEED_ON_START=true`. Seeds use `ON CONFLICT DO NOTHING`, so
restarts stay idempotent.

## LLM enrichment

Optional post-processing that runs **after** the base recap is built. It rewrites
the selected badges' copy and may prepend one `insight` scene right after the
`intro`. It is enabled only when `LLM_ENABLED=true` and `OPENAI_API_KEY` is set;
otherwise the base recap is served unchanged.

**Strictly additive / fail-open.** Enrichment never blocks or breaks the recap.
On timeout, transport error, provider error, or unparsable output, the original
payload is returned and the failure is logged. The base recap flow is untouched.

**Privacy: only public metrics leave the service.** The prompt is built from a
public allowlist derived from `metric_definitions.is_public`
([`RuleSet.PublicMetricKeys()`](internal/cards/rules.go)) — the single source of
truth shared with share cards. Private metrics (e.g. `moneyEarned`, `moneySaved`,
`sellerRating`, `messagesSent`, `activeListings`) are **never** sent to the
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
| `SEED_ON_START` | `false` | apply seeds from `SEEDS_DIR` on startup |
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

Requests from a listed origin get `Access-Control-Allow-Origin`; preflight
`OPTIONS` is answered with `204` and never reaches the routes. Requests from an
unlisted origin are still served, just without CORS headers, so the browser
blocks them. Responses always carry `Vary: Origin` so a shared cache cannot serve
one origin's headers to another.

Credentials are not enabled: no `Access-Control-Allow-Credentials`, so browsers
will not send cookies cross-origin. Add it here (and drop any wildcard) if the
API ever moves to cookie auth.
