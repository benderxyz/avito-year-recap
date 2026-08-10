# Architecture

The MVP is split into frontend, recap rendering packages, and three backend services.

## Services

- `frontend` - user journey: choose a test profile, generate recap, view cards, and continue to the next action.
- `recap-engine` - reusable packages for declarative recap scenes and React rendering.
- `analytics-service` - ingests user events and aggregates them into metrics, shares, and percentiles.
- `cards-service` - connects frontend with analytics and user data, builds recap cards from metrics: badges, insights, and recommendations, with optional LLM enrichment.
- `user-service` - manages test users and profiles.

## Data

The MVP uses ClickHouse for events (`analytics-service`) and Postgres for user profiles
(`user-service`, database `users`) and recap rules (`cards-service`, database `cards`).
Both Postgres databases live in a single instance and are created by `infra/postgres/init/`.
Each service owns its schema under `migrations/`.
Test profiles and seed events live in `seed-data/`. Sensitive fields, private messages, exact purchases, and other users' private data should not be exposed in personal recap cards or public share cards.
