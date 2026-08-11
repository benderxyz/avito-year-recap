DELETE FROM story_rules
WHERE id = 'stat-rating';

UPDATE story_rules
SET payload = '{"id":"insight-first-listing","type":"insight","title":"Первое объявление","text":"Первое объявление в этом году вы опубликовали {{value}}"}'::jsonb
WHERE id = 'insight-first-listing';
