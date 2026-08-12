UPDATE story_rules
SET payload = payload || '{"comparisonTemplate":"Объявлений у вас больше, чем у {{percentile}}% продавцов"}'::jsonb,
    updated_at = now()
WHERE id = 'stat-listings';

UPDATE story_rules
SET payload = payload || '{"comparisonTemplate":"Ваши объявления смотрели чаще, чем у {{percentile}}% продавцов"}'::jsonb,
    updated_at = now()
WHERE id = 'stat-views';

UPDATE story_rules
SET payload = payload || '{"comparisonTemplate":"В избранное вас добавляли чаще, чем у {{percentile}}% продавцов"}'::jsonb,
    updated_at = now()
WHERE id = 'stat-favorites';

UPDATE story_rules
SET payload = payload || '{"comparisonTemplate":"Вы переписывались активнее, чем {{percentile}}% пользователей"}'::jsonb,
    updated_at = now()
WHERE id = 'stat-messages';

UPDATE story_rules
SET payload = payload || '{"comparisonTemplate":"Сделок больше, чем у {{percentile}}% продавцов"}'::jsonb,
    updated_at = now()
WHERE id = 'stat-deals';
