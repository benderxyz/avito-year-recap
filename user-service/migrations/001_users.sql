CREATE TABLE IF NOT EXISTS users
(
    user_id      UInt64,
    external_id  String DEFAULT '',
    username     String DEFAULT '',
    timezone     String DEFAULT 'UTC',
    created_at   DateTime64(3),
    updated_at   DateTime64(3) DEFAULT now64(3)
)
ENGINE = ReplacingMergeTree(updated_at)
ORDER BY user_id
