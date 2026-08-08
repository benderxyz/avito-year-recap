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

INSERT INTO recommendation_rules (id, title, text, callout, link_label, path, priority, condition) VALUES
    ('resume-listings', 'Продвиньте свои объявления', 'У вас есть объявления, но сделок пока немного. Продвижение поможет им найти покупателя быстрее.', 'Больше показов — больше сделок', 'Продвинуть объявления', '/profile', 100, '{"match":"all","predicates":[{"metric":"listingsPublished","op":"gt","value":0},{"metric":"dealsClosed","op":"eq","value":0}]}'::jsonb),
    ('listings-favorited', 'Ваши объявления в избранном', 'Покупатели добавляли ваши объявления в избранное. Освежите их и ответьте на вопросы, пока интерес высокий.', 'Вами интересуются покупатели', 'К моим объявлениям', '/profile', 90, '{"predicates":[{"metric":"favoritesReceived","op":"gt","value":0}]}'::jsonb),
    ('resume-delivery', 'Авито Доставка', 'Вы уже пользовались доставкой. С ней проще продавать и покупать по всей стране.', 'Продавайте покупателям из других городов', 'Подключить доставку', '/dostavka', 80, '{"predicates":[{"metric":"deliveryOrders","op":"gt","value":0}]}'::jsonb),
    ('resume-active', 'Обновите активные объявления', 'У вас есть активные объявления. Освежите их, чтобы они снова поднялись в поиске.', 'Актуальные объявления смотрят чаще', 'К моим объявлениям', '/profile', 70, '{"predicates":[{"metric":"activeListings","op":"gt","value":0}]}'::jsonb),
    ('resume-search', 'Продолжите поиск', 'Вы активно искали в этом году. Возможно, нужное объявление уже появилось.', 'Новые объявления каждый день', 'Продолжить поиск', '/', 60, '{"match":"any","predicates":[{"metric":"searchQueries","op":"gt","value":0},{"metric":"viewsTotal","op":"gt","value":0}]}'::jsonb),
    ('new-listing', 'Разместите первое объявление', 'В этом году вы ещё не продавали. Начните — это займёт пару минут.', 'Продайте то, чем больше не пользуетесь', 'Разместить объявление', '/additem', 50, '{"predicates":[{"metric":"listingsPublished","op":"eq","value":0}]}'::jsonb)
ON CONFLICT (id) DO NOTHING;
