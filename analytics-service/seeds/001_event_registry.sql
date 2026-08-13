INSERT INTO event_registry (event_type, category, metric_key, unique_mode, unique_field, sort_order) VALUES
    ('item_published', 'counter', 'listingsPublished', NULL, NULL, 10),
    ('item_view_received', 'counter', 'viewsTotal', NULL, NULL, 20),
    ('item_favorited', 'counter', 'favoritesReceived', NULL, NULL, 30),
    ('chat_message_sent', 'counter', 'messagesSent', NULL, NULL, 40),
    ('deal_completed', 'counter', 'dealsClosed', NULL, NULL, 50),
    ('sale_proceeds', 'counter', 'moneyEarned', NULL, NULL, 60),
    ('buyer_saving', 'counter', 'moneySaved', NULL, NULL, 70),
    ('delivery_order_created', 'counter', 'deliveryOrders', NULL, NULL, 80),
    ('search_performed', 'counter', 'searchQueries', NULL, NULL, 90),
    ('active_items_count', 'gauge', 'activeListings', NULL, NULL, 100),
    ('peak_day_views', 'gauge', 'peakDayViews', NULL, NULL, 120),
    ('first_item_published', 'milestone', 'firstListingAt', NULL, NULL, 130),
    ('first_deal_completed', 'milestone', 'firstDealAt', NULL, NULL, 140),
    ('category_opened', 'unique', 'categoriesTried', 'payload', 'category', 150),
    ('day_active', 'unique', 'daysActive', 'day', NULL, 160),
    ('chat_reply', 'interval', 'avgReplySeconds', NULL, NULL, 170)
ON CONFLICT (event_type) DO UPDATE SET
    category = EXCLUDED.category,
    metric_key = EXCLUDED.metric_key,
    unique_mode = EXCLUDED.unique_mode,
    unique_field = EXCLUDED.unique_field,
    sort_order = EXCLUDED.sort_order,
    updated_at = now();

DELETE FROM event_registry
WHERE event_type NOT IN (
    'item_published', 'item_view_received', 'item_favorited', 'chat_message_sent',
    'deal_completed', 'sale_proceeds', 'buyer_saving', 'delivery_order_created',
    'search_performed', 'active_items_count', 'peak_day_views',
    'first_item_published', 'first_deal_completed', 'category_opened',
    'day_active', 'chat_reply'
);
