# Seed Data

Test users and a generator for Avito-like activity events.

## Users

`users/alex.json` — body for manual upsert if needed. The seed script also upserts a test user before posting events.

```sh
curl -s -X PUT http://localhost:8082/users/42 \
  -H 'Content-Type: application/json' \
  --data @seed-data/users/alex.json
```

## Events

Generate a year of activity and push it to the services

```sh
docker compose up -d clickhouse analytics user
go -C seed-data/seed-script run . -user 42 -year 2026
```

Flags

- `-user` user id (default 42)
- `-year` recap year (default 2026)
- `-analytics` analytics base url
- `-users` user-service base url
- `-seed` RNG seed for reproducible data

Then check metrics

```sh
curl -s 'http://localhost:8080/users/42/metrics?from=2026-01-01T00:00:00Z&to=2027-01-01T00:00:00Z'
```
