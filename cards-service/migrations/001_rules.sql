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

INSERT INTO badge_rules (id, title, description, icon_url, visibility, metric_key, op, threshold, sort_order) VALUES
    ('messages_sent_1000', 'Разговорчивый', 'Вы отправили больше 1000 сообщений!', '/badges/messages_sent_1000.svg', 'both', 'messagesSent', 'gt', 1000, 10),
    ('active_user', 'Активный пользователь', 'Вы заходили на нашу площадку в этом году!', '/badges/active_user.svg', 'both', 'daysActive', 'gt', 0, 20)
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_rules (id, scene_type, visibility, metric_key, op, threshold, payload, sort_order) VALUES
    ('stat-listings', 'stat', 'both', 'listingsPublished', 'gt', 0, '{"id":"stat-listings","type":"stat","value":"listingsPublished","unit":{"one":"объявление","few":"объявления","many":"объявлений"},"title":"вы опубликовали","eyebrow":"За год"}'::jsonb, 10),
    ('stat-views', 'stat', 'both', 'viewsTotal', 'gt', 0, '{"id":"stat-views","type":"stat","value":"viewsTotal","unit":{"one":"просмотр","few":"просмотра","many":"просмотров"},"title":"собрали ваши объявления"}'::jsonb, 20),
    ('stat-favorites', 'stat', 'both', 'favoritesReceived', 'gt', 0, '{"id":"stat-favorites","type":"stat","value":"favoritesReceived","unit":{"one":"добавление","few":"добавления","many":"добавлений"},"title":"в избранное","eyebrow":"Любимчики покупателей"}'::jsonb, 30),
    ('stat-messages', 'stat', 'private', 'messagesSent', 'gt', 0, '{"id":"stat-messages","type":"stat","value":"messagesSent","unit":{"one":"сообщение","few":"сообщения","many":"сообщений"},"title":"в чатах с покупателями","eyebrow":"Диалоги"}'::jsonb, 40),
    ('stat-deals', 'stat', 'both', 'dealsClosed', 'gt', 0, '{"id":"stat-deals","type":"stat","value":"dealsClosed","unit":{"one":"сделка","few":"сделки","many":"сделок"},"title":"успешно закрыто","eyebrow":"Результат"}'::jsonb, 50),
    ('stat-earned', 'stat', 'private', 'moneyEarned', 'gt', 0, '{"id":"stat-earned","type":"stat","value":"moneyEarned","unit":"₽","title":"заработали на продажах","eyebrow":"Доход","valueFormat":{"maximumFractionDigits":0}}'::jsonb, 60)
ON CONFLICT (id) DO NOTHING;
