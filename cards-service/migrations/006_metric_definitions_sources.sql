ALTER TABLE metric_definitions ADD COLUMN IF NOT EXISTS source_key text;

ALTER TABLE metric_definitions ADD COLUMN IF NOT EXISTS source_field text NOT NULL DEFAULT 'value';

ALTER TABLE metric_definitions ADD COLUMN IF NOT EXISTS include_in_llm boolean NOT NULL DEFAULT false;

ALTER TABLE metric_definitions DROP CONSTRAINT IF EXISTS metric_definitions_value_type_check;

ALTER TABLE metric_definitions ADD CONSTRAINT metric_definitions_value_type_check
    CHECK (value_type IN ('number', 'money', 'percentile', 'ratio', 'string', 'date'));

ALTER TABLE metric_definitions DROP CONSTRAINT IF EXISTS metric_definitions_source_field_check;

ALTER TABLE metric_definitions ADD CONSTRAINT metric_definitions_source_field_check
    CHECK (source_field IN ('value', 'percentile', 'share'));

UPDATE metric_definitions SET source_key = key WHERE source_key IS NULL;

UPDATE metric_definitions AS d
SET source_key = p.key,
    source_field = 'percentile'
FROM metric_definitions AS p
WHERE p.percentile_key = d.key
  AND d.value_type = 'percentile';

UPDATE metric_definitions
SET value_type = 'date'
WHERE key IN ('firstListingAt', 'firstDealAt');

UPDATE metric_definitions
SET include_in_llm = true
WHERE key IN (
    'listingsPublished', 'viewsTotal', 'favoritesReceived', 'messagesSent',
    'dealsClosed', 'moneyEarned', 'moneySaved', 'daysActive', 'peakDayViews',
    'searchQueries', 'categoriesTried', 'deliveryOrders', 'activeListings',
    'sellerRating'
);
