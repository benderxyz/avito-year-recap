CREATE TABLE IF NOT EXISTS users (
    user_id     bigint PRIMARY KEY,
    external_id text NOT NULL DEFAULT '',
    username    text NOT NULL DEFAULT '',
    timezone    text NOT NULL DEFAULT 'UTC',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);
