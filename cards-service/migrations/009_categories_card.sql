UPDATE story_rules
SET scene_type = 'blocks',
    payload = '{"id":"stat-categories","type":"blocks","blocks":[{"type":"stat","value":"categoriesTried","unit":{"one":"категория","few":"категории","many":"категорий"},"title":"вы открывали за год","eyebrow":"Кругозор"},{"type":"callout","text":"Заглядывали в разные разделы, от электроники до садовой мебели"}]}'::jsonb,
    updated_at = now()
WHERE id = 'stat-categories';
