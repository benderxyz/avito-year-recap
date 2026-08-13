ALTER TABLE metric_definitions
    ADD COLUMN IF NOT EXISTS comparison_min_percentile numeric;

ALTER TABLE metric_definitions
    DROP CONSTRAINT IF EXISTS metric_definitions_comparison_min_percentile_check;

ALTER TABLE metric_definitions ADD CONSTRAINT metric_definitions_comparison_min_percentile_check
    CHECK (comparison_min_percentile IS NULL
        OR (comparison_min_percentile >= 0 AND comparison_min_percentile <= 100));
