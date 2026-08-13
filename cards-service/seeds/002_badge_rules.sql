INSERT INTO badge_rules (id, title, description, icon_url, visibility, metric_key, op, threshold, sort_order) VALUES
    ('messages_sent_1000', 'Разговорчивый', 'Вы отправили больше 1000 сообщений!', '/badges/messages_sent_1000.svg', 'both', 'messagesSent', 'gt', 1000, 10),
    ('active_user', 'Активный пользователь', 'Вы заходили на нашу площадку в этом году!', '/badges/active_user.svg', 'both', 'daysActive', 'gt', 0, 20)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    icon_url = EXCLUDED.icon_url,
    visibility = EXCLUDED.visibility,
    metric_key = EXCLUDED.metric_key,
    op = EXCLUDED.op,
    threshold = EXCLUDED.threshold,
    sort_order = EXCLUDED.sort_order,
    updated_at = now();
