# Analytics Service

Service responsible for user event storage, aggregation, metrics, achievements, and reproducible recap rules.

Ownership: Maxim Shugedov.

Planned areas:

- `internal/events` - raw user activity events.
- `internal/aggregation` - calculated metrics and statistics.
- `internal/rules` - reproducible MVP rules for insights, badges, and recommendations.
- `internal/db` - database access.
- `internal/api` - service API handlers.
