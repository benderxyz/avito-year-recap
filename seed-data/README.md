# Seed Data

Test users plus a generator for Avito-like activity events.

## Users

`users/` has JSON bodies for manual upsert if needed. The seed script upserts the whole roster before posting events.

Default roster (ids 42-51) has different activity levels. That way global share plus percentile in analytics have something to compare against.

```sh
curl -s -X PUT http://localhost:8082/users/42 \
  -H 'Content-Type: application/json' \
  --data @seed-data/users/alex.json
```

## Events

Generate a year of activity then push it to the services.

```sh
docker compose up -d postgres clickhouse user analytics
go -C seed-data/seed-script run . -year 2026
```

Postgres is one container. user-service writes profiles into database `users`. Wait until `curl localhost:8082/health` and `curl localhost:8080/health` succeed before running the script.

By default this seeds 10 users (42-51). Pass `-user` to limit the set.

```sh
go -C seed-data/seed-script run . -user 42,43,48 -year 2026
```

Flags

* `-user` comma-separated user ids. Empty means the default roster.
* `-year` recap year (default 2026)
* `-analytics` analytics base url
* `-users` user-service base url
* `-seed` RNG seed for reproducible data

Then check metrics.

```sh
curl -s 'http://localhost:8080/users/42/metrics?from=2026-01-01T00:00:00Z&to=2027-01-01T00:00:00Z'
curl -s 'http://localhost:8080/users/47/metrics?from=2026-01-01T00:00:00Z&to=2027-01-01T00:00:00Z'
```
