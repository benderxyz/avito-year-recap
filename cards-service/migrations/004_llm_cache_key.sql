DROP TABLE IF EXISTS recap_llm_cache;

CREATE TABLE IF NOT EXISTS recap_llm_cache (
    user_id        text        NOT NULL,
    year           integer     NOT NULL,
    mode           text        NOT NULL,
    prompt_version text        NOT NULL,
    result         jsonb       NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, year, mode, prompt_version)
);
