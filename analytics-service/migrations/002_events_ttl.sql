ALTER TABLE events
    MODIFY TTL toDateTime(occurred_at) + INTERVAL 2 YEAR
