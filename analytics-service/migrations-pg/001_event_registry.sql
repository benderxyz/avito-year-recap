CREATE TABLE IF NOT EXISTS event_registry (
    event_type   text PRIMARY KEY,
    category     text NOT NULL CHECK (category IN ('counter', 'interval', 'gauge', 'milestone', 'unique')),
    metric_key   text NOT NULL,
    unique_mode  text CHECK (unique_mode IN ('payload', 'day')),
    unique_field text,
    enabled      boolean NOT NULL DEFAULT true,
    sort_order   integer NOT NULL DEFAULT 0,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);
