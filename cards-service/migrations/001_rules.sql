CREATE TABLE IF NOT EXISTS badge_rules (
    id text PRIMARY KEY,
    title text NOT NULL,
    description text NOT NULL,
    icon_url text,
    visibility text NOT NULL DEFAULT 'both' CHECK (visibility IN ('private', 'public', 'both')),
    metric_key text NOT NULL,
    op text NOT NULL CHECK (op IN ('gt', 'gte', 'eq', 'exists')),
    threshold double precision NOT NULL DEFAULT 0,
    sort_order integer NOT NULL DEFAULT 0,
    enabled boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE badge_rules ADD COLUMN IF NOT EXISTS icon_url text;

CREATE TABLE IF NOT EXISTS story_rules (
    id text PRIMARY KEY,
    scene_type text NOT NULL,
    visibility text NOT NULL DEFAULT 'both' CHECK (visibility IN ('private', 'public', 'both')),
    metric_key text,
    op text CHECK (op IN ('gt', 'gte', 'eq', 'exists')),
    threshold double precision,
    payload jsonb NOT NULL,
    sort_order integer NOT NULL DEFAULT 0,
    enabled boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS recommendation_rules (
    id text PRIMARY KEY,
    title text NOT NULL,
    text text NOT NULL,
    callout text NOT NULL DEFAULT '',
    link_label text NOT NULL,
    path text NOT NULL,
    priority integer NOT NULL DEFAULT 0,
    condition jsonb NOT NULL DEFAULT '{}'::jsonb,
    enabled boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
