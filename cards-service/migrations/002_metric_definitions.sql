CREATE TABLE IF NOT EXISTS metric_definitions (
    key text PRIMARY KEY,
    value_type text NOT NULL CHECK (value_type IN ('number', 'money', 'percentile', 'string')),
    currency text,
    is_public boolean NOT NULL DEFAULT false,
    percentile_key text,
    enabled boolean NOT NULL DEFAULT true,
    sort_order integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
