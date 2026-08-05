CREATE TABLE IF NOT EXISTS events
(
    event_id       UUID DEFAULT generateUUIDv4(),
    user_id        UInt64,
    session_id     UUID,

    event_type     LowCardinality(String),
    event_category Enum8(
        'counter'   = 1,
        'interval'  = 2,
        'gauge'     = 3,
        'milestone' = 4,
        'unique'    = 5
    ),

    value          Float64 DEFAULT 0,
    payload        String DEFAULT '{}',

    occurred_at    DateTime64(3) DEFAULT now64(3),
    inserted_at    DateTime64(3) DEFAULT now64(3)
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(occurred_at)
ORDER BY (user_id, event_type, occurred_at)
TTL toDateTime(occurred_at) + INTERVAL 2 YEAR
