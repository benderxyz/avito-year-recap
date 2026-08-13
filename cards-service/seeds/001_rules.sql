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

INSERT INTO story_rules (id, scene_type, visibility, metric_key, op, threshold, payload, sort_order) VALUES
    ('intro', 'intro', 'both', NULL, NULL, NULL, '{"id":"intro","type":"intro","title":"{{displayName}}, ваш {{year}} на Авито","subtitle":"Год находок, сделок и удачных объявлений","actions":[{"type":"next","label":"Смотреть итоги"}]}'::jsonb, 1),
    ('stat-listings', 'stat', 'both', 'listingsPublished', 'gt', 0, '{"id":"stat-listings","type":"stat","value":"listingsPublished","unit":{"one":"объявление","few":"объявления","many":"объявлений"},"title":"вы опубликовали","eyebrow":"За год","comparisonTemplate":"Объявлений у вас больше, чем у {{percentile}}% продавцов"}'::jsonb, 10),
    ('stat-views', 'stat', 'both', 'viewsTotal', 'gt', 0, '{"id":"stat-views","type":"stat","value":"viewsTotal","unit":{"one":"просмотр","few":"просмотра","many":"просмотров"},"title":"собрали ваши объявления","comparisonTemplate":"Ваши объявления смотрели чаще, чем у {{percentile}}% продавцов"}'::jsonb, 20),
    ('stat-favorites', 'stat', 'both', 'favoritesReceived', 'gt', 0, '{"id":"stat-favorites","type":"stat","value":"favoritesReceived","unit":{"one":"добавление","few":"добавления","many":"добавлений"},"title":"в избранное","eyebrow":"Любимчики покупателей","comparisonTemplate":"В избранное вас добавляли чаще, чем у {{percentile}}% продавцов"}'::jsonb, 30),
    ('stat-messages', 'stat', 'private', 'messagesSent', 'gt', 0, '{"id":"stat-messages","type":"stat","value":"messagesSent","unit":{"one":"сообщение","few":"сообщения","many":"сообщений"},"title":"в чатах с покупателями","eyebrow":"Диалоги","comparisonTemplate":"Вы переписывались активнее, чем {{percentile}}% пользователей"}'::jsonb, 40),
    ('stat-deals', 'stat', 'both', 'dealsClosed', 'gt', 0, '{"id":"stat-deals","type":"stat","value":"dealsClosed","unit":{"one":"сделка","few":"сделки","many":"сделок"},"title":"успешно закрыто","eyebrow":"Результат","comparisonTemplate":"Сделок больше, чем у {{percentile}}% продавцов"}'::jsonb, 50),
    ('stat-earned', 'stat', 'private', 'moneyEarned', 'gt', 0, '{"id":"stat-earned","type":"stat","value":"moneyEarned","unit":"₽","title":"заработали на продажах","eyebrow":"Доход","valueFormat":{"maximumFractionDigits":0}}'::jsonb, 60),
    ('stat-saved', 'stat', 'private', 'moneySaved', 'gt', 0, '{"id":"stat-saved","type":"stat","value":"moneySaved","unit":"₽","title":"сэкономили на покупках","eyebrow":"Выгода","valueFormat":{"maximumFractionDigits":0}}'::jsonb, 70),
    ('blocks-days-active', 'blocks', 'both', 'daysActive', 'gt', 0, '{"id":"blocks-days-active","type":"blocks","blocks":[{"type":"stat","value":"daysActive","unit":{"one":"день","few":"дня","many":"дней"},"title":"были активны","eyebrow":"Присутствие"},{"type":"callout","text":"Почти каждый день на Авито что-то происходило"}]}'::jsonb, 80),
    ('stat-peak-views', 'stat', 'both', 'peakDayViews', 'gt', 0, '{"id":"stat-peak-views","type":"stat","value":"peakDayViews","unit":{"one":"просмотр","few":"просмотра","many":"просмотров"},"title":"за один день","eyebrow":"Рекорд"}'::jsonb, 90),
    ('stat-categories', 'blocks', 'both', 'categoriesTried', 'gt', 0, '{"id":"stat-categories","type":"blocks","blocks":[{"type":"stat","value":"categoriesTried","unit":{"one":"категория","few":"категории","many":"категорий"},"title":"вы открывали за год","eyebrow":"Кругозор"},{"type":"callout","text":"Заглядывали в разные разделы, от электроники до садовой мебели"}]}'::jsonb, 100),
    ('blocks-search', 'blocks', 'both', 'searchQueries', 'gt', 0, '{"id":"blocks-search","type":"blocks","blocks":[{"type":"stat","value":"searchQueries","unit":{"one":"поиск","few":"поиска","many":"поисков"},"title":"запросов за год","eyebrow":"Любопытство"},{"type":"callout","text":"Вы не только продавали — вы и сами охотились за находками"}]}'::jsonb, 110),
    ('stat-delivery', 'stat', 'both', 'deliveryOrders', 'gt', 0, '{"id":"stat-delivery","type":"stat","value":"deliveryOrders","unit":{"one":"заказ","few":"заказа","many":"заказов"},"title":"через Авито Доставку","eyebrow":"Доставка"}'::jsonb, 120),
    ('stat-active-listings', 'stat', 'private', 'activeListings', 'gt', 0, '{"id":"stat-active-listings","type":"stat","value":"activeListings","unit":{"one":"объявление","few":"объявления","many":"объявлений"},"title":"сейчас активны","eyebrow":"В продаже"}'::jsonb, 130),
    ('stat-reply', 'stat', 'private', 'avgReplySeconds', 'gt', 0, '{"id":"stat-reply","type":"stat","value":"avgReplySeconds","unit":{"one":"секунда","few":"секунды","many":"секунд"},"title":"в среднем до ответа","eyebrow":"Скорость","valueFormat":{"maximumFractionDigits":0}}'::jsonb, 150),
    ('insight-first-listing', 'insight', 'private', 'firstListingAt', 'gt', 0, '{"id":"insight-first-listing","type":"insight","value":"firstListingAt","title":"Первое объявление","text":"Первое объявление в этом году вы опубликовали {{value}}"}'::jsonb, 160),
    ('insight-first-deal', 'insight', 'private', 'firstDealAt', 'gt', 0, '{"id":"insight-first-deal","type":"insight","value":"firstDealAt","title":"Первая сделка","text":"Первую сделку закрыли {{value}}"}'::jsonb, 170),
    ('outro', 'outro', 'both', NULL, NULL, NULL, '{"id":"outro","type":"outro","title":"Это был ваш год на Авито","subtitle":"Сохраните итоги или вернитесь к объявлениям","actions":[{"type":"share","label":"Поделиться","share":{"kind":"link","title":"Мои итоги на Авито","text":"Посмотрите, каким был мой год на Авито!"}},{"type":"custom","id":"close-recap","label":"На главную","variant":"primary"}]}'::jsonb, 1000)
ON CONFLICT (id) DO UPDATE SET
    scene_type = EXCLUDED.scene_type,
    visibility = EXCLUDED.visibility,
    metric_key = EXCLUDED.metric_key,
    op = EXCLUDED.op,
    threshold = EXCLUDED.threshold,
    payload = EXCLUDED.payload,
    sort_order = EXCLUDED.sort_order,
    updated_at = now();

INSERT INTO recommendation_rules (id, title, text, callout, link_label, path, priority, condition) VALUES
    ('resume-listings', 'Продвиньте свои объявления', 'У вас есть объявления, но сделок пока немного. Продвижение поможет им найти покупателя быстрее.', 'Больше показов — больше сделок', 'Продвинуть объявления', '/profile', 100, '{"match":"all","predicates":[{"metric":"listingsPublished","op":"gt","value":0},{"metric":"dealsClosed","op":"eq","value":0}]}'::jsonb),
    ('listings-favorited', 'Ваши объявления в избранном', 'Покупатели добавляли ваши объявления в избранное. Освежите их и ответьте на вопросы, пока интерес высокий.', 'Вами интересуются покупатели', 'К моим объявлениям', '/profile', 90, '{"predicates":[{"metric":"favoritesReceived","op":"gt","value":0}]}'::jsonb),
    ('resume-delivery', 'Авито Доставка', 'Вы уже пользовались доставкой. С ней проще продавать и покупать по всей стране.', 'Продавайте покупателям из других городов', 'Подключить доставку', '/dostavka', 80, '{"predicates":[{"metric":"deliveryOrders","op":"gt","value":0}]}'::jsonb),
    ('resume-active', 'Обновите активные объявления', 'У вас есть активные объявления. Освежите их, чтобы они снова поднялись в поиске.', 'Актуальные объявления смотрят чаще', 'К моим объявлениям', '/profile', 70, '{"predicates":[{"metric":"activeListings","op":"gt","value":0}]}'::jsonb),
    ('resume-search', 'Продолжите поиск', 'Вы активно искали в этом году. Возможно, нужное объявление уже появилось.', 'Новые объявления каждый день', 'Продолжить поиск', '/', 60, '{"match":"any","predicates":[{"metric":"searchQueries","op":"gt","value":0},{"metric":"viewsTotal","op":"gt","value":0}]}'::jsonb),
    ('new-listing', 'Разместите первое объявление', 'В этом году вы ещё не продавали. Начните — это займёт пару минут.', 'Продайте то, чем больше не пользуетесь', 'Разместить объявление', '/additem', 50, '{"predicates":[{"metric":"listingsPublished","op":"eq","value":0}]}'::jsonb)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    text = EXCLUDED.text,
    callout = EXCLUDED.callout,
    link_label = EXCLUDED.link_label,
    path = EXCLUDED.path,
    priority = EXCLUDED.priority,
    condition = EXCLUDED.condition,
    updated_at = now();

DELETE FROM badge_rules
WHERE id NOT IN ('messages_sent_1000', 'active_user');

DELETE FROM story_rules
WHERE id NOT IN (
    'intro', 'stat-listings', 'stat-views', 'stat-favorites', 'stat-messages',
    'stat-deals', 'stat-earned', 'stat-saved', 'blocks-days-active',
    'stat-peak-views', 'stat-categories', 'blocks-search', 'stat-delivery',
    'stat-active-listings', 'stat-reply', 'insight-first-listing',
    'insight-first-deal', 'outro'
);

DELETE FROM recommendation_rules
WHERE id NOT IN (
    'resume-listings', 'listings-favorited', 'resume-delivery', 'resume-active',
    'resume-search', 'new-listing'
);
