# Architecture

The MVP is split into frontend, recap rendering packages, and three backend services.

## Services

- `frontend` - user journey: choose a test profile, generate recap, view cards, and continue to the next action.
- `recap-engine` - reusable packages for declarative recap scenes and React rendering.
- `analytics-service` - stores and aggregates user events, calculates metrics, badges, insights, and recommendations.
- `cards-service` - connects frontend with analytics and user data, prepares recap cards, and integrates with neural network generation.
- `user-service` - manages test users, profiles, and user-related event intake.

## Data

The MVP uses ClickHouse for events (`analytics-service`) and user profiles (`user-service`).
Test profiles and seed events live in `seed-data/`. Sensitive fields, private messages, exact purchases, and other users' private data should not be exposed in personal recap cards or public share cards.

## Future Work

- Define service interaction protocol in `contracts`.
- Add stubs for frontend-to-service and service-to-service communication.
- Document database choice and migrations.
