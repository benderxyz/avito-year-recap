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
