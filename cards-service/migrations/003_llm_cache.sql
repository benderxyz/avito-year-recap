CREATE TABLE IF NOT EXISTS recap_llm_cache (
    user_id    text        NOT NULL,
    year       integer     NOT NULL,
    result     jsonb       NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, year)
);
