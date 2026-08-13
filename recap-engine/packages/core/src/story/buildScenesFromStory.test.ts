import { describe, expect, it } from '@jest/globals';

import { createFormatters } from '../format/createFormatters';
import { resolveValue } from '../resolve/resolveValue';
import { createTheme } from '../theme/createTheme';
import { ESceneActionType } from '../types/actions';
import type { RecapContext, ValueOrFn } from '../types/context';
import { EMotionPreset } from '../types/motion';
import { EMetricType, type RecapPayload } from '../types/payload';
import { ESceneBlockType, ESceneType, type SceneDefinition } from '../types/scenes';
import { buildScenesFromStory } from './buildScenesFromStory';

const makePayload = (overrides: Partial<RecapPayload> = {}): RecapPayload => ({
  schemaVersion: 1,
  meta: {
    vertical: 'general',
    year: 2025,
    locale: 'ru-RU',
    user: { id: 'user-1', displayName: 'Максим' },
    generatedAt: '2026-01-01T00:00:00.000Z',
  },
  metrics: {
    sales: { type: EMetricType.Number, value: 42 },
    rank: { type: EMetricType.Percentile, value: 91 },
    saving: { type: EMetricType.Money, value: 1250, currency: 'RUB' },
  },
  story: [],
  ...overrides,
});

const allTypesPayload = makePayload({
  narrative: {
    scenes: {
      intro: { title: 'Narrative intro', subtitle: 'Narrative subtitle' },
      stat: { title: 'Narrative stat', comparison: 'топ {{percentile}}' },
      insight: { title: 'Narrative insight', body: 'Narrative body' },
      upsell: {
        title: 'Narrative upsell',
        body: 'Экономия {{value}} и снова {{value}}',
        comparison: 'Narrative callout {{value}}',
      },
      outro: { title: 'Narrative outro', subtitle: 'Narrative goodbye' },
    },
  },
  story: [
    {
      id: 'intro',
      type: ESceneType.Intro,
      title: 'Item intro',
      subtitle: 'Item subtitle',
      motion: EMotionPreset.Fade,
      durationMs: 700,
    },
    {
      id: 'stat',
      type: ESceneType.Stat,
      value: 'sales',
      percentile: 'rank',
      unit: { one: 'продажа', few: 'продажи', many: 'продаж' },
      title: 'Item stat',
      valueFormat: { maximumFractionDigits: 1 },
    },
    {
      id: 'insight',
      type: ESceneType.Insight,
      title: 'Item insight',
      text: 'Item body',
      linksTo: 'stat',
    },
    {
      id: 'achievement',
      type: ESceneType.Achievement,
      badgeId: 'seller',
      title: 'Продавец года',
      description: '42 продажи',
      icon: '/badge.svg',
    },
    {
      id: 'upsell',
      type: ESceneType.Upsell,
      title: 'Item upsell',
      text: 'Item body {{value}}',
      callout: 'Item callout {{value}}',
      value: 'saving',
    },
    {
      id: 'blocks',
      type: ESceneType.Blocks,
      actions: [{ type: ESceneActionType.Next, label: 'Дальше' }],
      blocks: [
        {
          type: ESceneBlockType.Stat,
          value: 'sales',
          percentile: 'rank',
          title: 'Block stat',
          eyebrow: 'Block eyebrow',
          unit: 'шт.',
          valueFormat: { notation: 'compact' },
          blockMotion: EMotionPreset.None,
        },
        {
          type: ESceneBlockType.Text,
          text: 'Block text',
          blockMotion: EMotionPreset.None,
        },
        {
          type: ESceneBlockType.Callout,
          text: 'Block callout',
          blockMotion: EMotionPreset.None,
        },
      ],
    },
    {
      id: 'outro',
      type: ESceneType.Outro,
      title: 'Item outro',
      subtitle: 'Item goodbye',
    },
    {
      id: 'custom',
      type: ESceneType.Custom,
      sceneType: 'timeline',
      props: { highlight: true },
    },
  ],
});

const contextFor = (payload: RecapPayload): RecapContext<RecapPayload> => ({
  data: payload,
  theme: createTheme(),
  index: 0,
  total: payload.story.length,
  format: createFormatters('ru-RU'),
});

function resolved<T>(
  value: ValueOrFn<T, RecapPayload> | undefined,
  payload = allTypesPayload,
): T | undefined {
  return resolveValue(value, contextFor(payload));
}

function sceneById(
  scenes: SceneDefinition<RecapPayload>[],
  id: string,
): SceneDefinition<RecapPayload> {
  const scene = scenes.find((candidate) => candidate.id === id);
  if (!scene) throw new Error(`Missing test scene: ${id}`);
  return scene;
}

describe('buildScenesFromStory', () => {
  const scenes = buildScenesFromStory(allTypesPayload);

  it('maps every scene type in playlist order', () => {
    expect(scenes.map(({ id, type }) => [id, type])).toEqual([
      ['intro', ESceneType.Intro],
      ['stat', ESceneType.Stat],
      ['insight', ESceneType.Insight],
      ['achievement', ESceneType.Achievement],
      ['upsell', ESceneType.Upsell],
      ['blocks', ESceneType.Blocks],
      ['outro', ESceneType.Outro],
      ['custom', ESceneType.Custom],
    ]);
  });

  it('preserves common configuration and supplies intro actions by default', () => {
    const intro = sceneById(scenes, 'intro');
    expect(intro).toMatchObject({
      motion: EMotionPreset.Fade,
      durationMs: 700,
      actions: [{ type: ESceneActionType.Next, label: 'Начать' }],
    });

    const blocks = sceneById(scenes, 'blocks');
    expect(blocks.actions).toEqual([{ type: ESceneActionType.Next, label: 'Дальше' }]);
  });

  it('uses narrative intro copy ahead of item copy', () => {
    const intro = sceneById(scenes, 'intro');
    expect(intro.type).toBe(ESceneType.Intro);
    if (intro.type !== ESceneType.Intro) return;

    expect(resolved(intro.title)).toBe('Narrative intro');
    expect(resolved(intro.subtitle)).toBe('Narrative subtitle');
  });

  it('maps stat metrics, formatting, defaults and narrative comparison', () => {
    const stat = sceneById(scenes, 'stat');
    expect(stat.type).toBe(ESceneType.Stat);
    if (stat.type !== ESceneType.Stat) return;

    expect(resolved(stat.eyebrow)).toBe('За 2025 год');
    expect(resolved(stat.title)).toBe('Narrative stat');
    expect(resolved(stat.value)).toBe(42);
    expect(stat.unit).toEqual({ one: 'продажа', few: 'продажи', many: 'продаж' });
    expect(stat.valueFormat).toEqual({ maximumFractionDigits: 1 });
    expect(stat.blockMotion).toBe(EMotionPreset.CountUp);
    expect(stat.comparison?.template).toBe('топ {{percentile}}');
    expect(resolved(stat.comparison?.percentile)).toBe(91);
  });

  it('maps insight narrative copy, link and default motion', () => {
    const insight = sceneById(scenes, 'insight');
    expect(insight.type).toBe(ESceneType.Insight);
    if (insight.type !== ESceneType.Insight) return;

    expect(resolved(insight.title)).toBe('Narrative insight');
    expect(resolved(insight.text)).toBe('Narrative body');
    expect(insight.linksTo).toBe('stat');
    expect(insight.blockMotion).toBe(EMotionPreset.StaggerText);
  });

  it('maps achievement fields and default badge motion', () => {
    expect(sceneById(scenes, 'achievement')).toMatchObject({
      type: ESceneType.Achievement,
      badgeId: 'seller',
      title: 'Продавец года',
      description: '42 продажи',
      icon: '/badge.svg',
      blockMotion: EMotionPreset.BadgePop,
    });
  });

  it('formats every money placeholder in upsell copy and prefers item callout', () => {
    const upsell = sceneById(scenes, 'upsell');
    expect(upsell.type).toBe(ESceneType.Upsell);
    if (upsell.type !== ESceneType.Upsell) return;

    const currency = contextFor(allTypesPayload).format.currency(1250);
    expect(resolved(upsell.title)).toBe('Narrative upsell');
    expect(resolved(upsell.text)).toBe(`Экономия ${currency} и снова ${currency}`);
    expect(resolved(upsell.callout)).toBe(`Item callout ${currency}`);
    expect(upsell.blockMotion).toBe(EMotionPreset.CalloutIn);
  });

  it.each(['firstListing', 'sales', 'saving', 'city'])(
    'fills the value placeholder in insight copy for the %s metric',
    (key) => {
      const payload = makePayload({
        metrics: {
          sales: { type: EMetricType.Number, value: 42 },
          saving: { type: EMetricType.Money, value: 1250, currency: 'RUB' },
          city: { type: EMetricType.String, value: 'Москва' },
          firstListing: { type: EMetricType.Date, value: '2024-03-14' },
        },
        story: [
          {
            id: 'insight',
            type: ESceneType.Insight,
            text: 'Первое объявление {{value}}',
            value: key,
          },
        ],
      });

      const insight = sceneById(buildScenesFromStory(payload), 'insight');
      if (insight.type !== ESceneType.Insight) return;

      const format = contextFor(payload).format;
      const expected: Record<string, string> = {
        firstListing: `Первое объявление ${format.date(new Date('2024-03-14'))}`,
        sales: `Первое объявление ${format.number(42)}`,
        saving: `Первое объявление ${format.currency(1250, 'RUB')}`,
        city: 'Первое объявление Москва',
      };

      expect(resolved(insight.text, payload)).toBe(expected[key]);
    },
  );

  it('keeps the value placeholder when the metric is missing', () => {
    const payload = makePayload({
      story: [
        {
          id: 'insight',
          type: ESceneType.Insight,
          text: 'Первое объявление {{value}}',
          value: 'unknownMetric',
        },
      ],
    });

    const insight = sceneById(buildScenesFromStory(payload), 'insight');
    if (insight.type !== ESceneType.Insight) return;

    expect(resolved(insight.text, payload)).toBe('Первое объявление {{value}}');
  });

  it('applies the requested date format to insight copy', () => {
    const payload = makePayload({
      metrics: { firstListing: { type: EMetricType.Date, value: '2024-03-14' } },
      story: [
        {
          id: 'insight',
          type: ESceneType.Insight,
          text: 'Начало {{value}}',
          value: 'firstListing',
          dateFormat: { day: 'numeric', month: 'numeric', year: 'numeric' },
        },
      ],
    });

    const insight = sceneById(buildScenesFromStory(payload), 'insight');
    if (insight.type !== ESceneType.Insight) return;

    expect(resolved(insight.text, payload)).toBe('Начало 14.03.2024');
  });

  it('maps stat, text and callout blocks with their options', () => {
    const blocksScene = sceneById(scenes, 'blocks');
    expect(blocksScene.type).toBe(ESceneType.Blocks);
    if (blocksScene.type !== ESceneType.Blocks) return;

    expect(blocksScene.blocks).toHaveLength(3);
    const [stat, text, callout] = blocksScene.blocks;
    expect(stat.type).toBe(ESceneBlockType.Stat);
    if (stat.type !== ESceneBlockType.Stat) return;
    expect(resolved(stat.value)).toBe(42);
    expect(resolved(stat.title)).toBe('Block stat');
    expect(resolved(stat.eyebrow)).toBe('Block eyebrow');
    expect(stat).toMatchObject({
      unit: 'шт.',
      valueFormat: { notation: 'compact' },
      blockMotion: EMotionPreset.None,
    });
    expect(stat.comparison?.template).toBe('больше, чем у {{percentile}}% пользователей');
    expect(resolved(stat.comparison?.percentile)).toBe(91);
    expect(text).toEqual({
      type: ESceneBlockType.Text,
      text: 'Block text',
      blockMotion: EMotionPreset.None,
    });
    expect(callout).toEqual({
      type: ESceneBlockType.Callout,
      text: 'Block callout',
      blockMotion: EMotionPreset.None,
    });
  });

  it('maps outro narrative copy and custom scene data', () => {
    const outro = sceneById(scenes, 'outro');
    expect(outro.type).toBe(ESceneType.Outro);
    if (outro.type !== ESceneType.Outro) return;
    expect(resolved(outro.title)).toBe('Narrative outro');
    expect(resolved(outro.subtitle)).toBe('Narrative goodbye');

    expect(sceneById(scenes, 'custom')).toMatchObject({
      type: ESceneType.Custom,
      sceneType: 'timeline',
      props: { highlight: true },
    });
  });

  it('uses item copy and explicit motion overrides when narrative is absent', () => {
    const payload = makePayload({
      story: [
        {
          id: 'stat',
          type: ESceneType.Stat,
          value: 'sales',
          title: 'Item title',
          eyebrow: '',
          blockMotion: EMotionPreset.None,
        },
        {
          id: 'insight',
          type: ESceneType.Insight,
          title: 'Insight title',
          text: 'Insight text',
          blockMotion: EMotionPreset.None,
        },
        {
          id: 'upsell',
          type: ESceneType.Upsell,
          title: 'Upsell title',
          text: 'No placeholder',
          callout: 'Callout',
          blockMotion: EMotionPreset.None,
        },
      ],
    });
    const mapped = buildScenesFromStory(payload);
    const stat = sceneById(mapped, 'stat');
    const insight = sceneById(mapped, 'insight');
    const upsell = sceneById(mapped, 'upsell');

    if (stat.type !== ESceneType.Stat || insight.type !== ESceneType.Insight) return;
    if (upsell.type !== ESceneType.Upsell) return;
    expect(resolved(stat.title, payload)).toBe('Item title');
    expect(resolved(stat.eyebrow, payload)).toBe('');
    expect(stat.comparison).toBeUndefined();
    expect(stat.blockMotion).toBe(EMotionPreset.None);
    expect(resolved(insight.title, payload)).toBe('Insight title');
    expect(resolved(insight.text, payload)).toBe('Insight text');
    expect(insight.blockMotion).toBe(EMotionPreset.None);
    expect(resolved(upsell.text, payload)).toBe('No placeholder');
    expect(resolved(upsell.callout, payload)).toBe('Callout');
    expect(upsell.blockMotion).toBe(EMotionPreset.None);
  });

  it('supplies copy and comparison defaults', () => {
    const payload = makePayload({
      story: [
        { id: 'intro', type: ESceneType.Intro },
        {
          id: 'stat',
          type: ESceneType.Stat,
          value: 'missing',
          percentile: 'missing-percentile',
        },
        { id: 'insight', type: ESceneType.Insight },
        { id: 'upsell', type: ESceneType.Upsell },
        { id: 'outro', type: ESceneType.Outro },
      ],
    });
    const mapped = buildScenesFromStory(payload);
    const intro = sceneById(mapped, 'intro');
    const stat = sceneById(mapped, 'stat');
    const insight = sceneById(mapped, 'insight');
    const upsell = sceneById(mapped, 'upsell');
    const outro = sceneById(mapped, 'outro');

    if (intro.type !== ESceneType.Intro || stat.type !== ESceneType.Stat) return;
    if (insight.type !== ESceneType.Insight || upsell.type !== ESceneType.Upsell) return;
    if (outro.type !== ESceneType.Outro) return;
    expect(resolved(intro.title, payload)).toBe('Максим, ваш 2025');
    expect(resolved(intro.subtitle, payload)).toBe('Год на Авито в цифрах');
    expect(resolved(stat.title, payload)).toBe('');
    expect(resolved(stat.value, payload)).toBe(0);
    expect(stat.comparison?.template).toBe('это больше, чем у {{percentile}}% пользователей');
    expect(resolved(stat.comparison?.percentile, payload)).toBe(0);
    expect(resolved(insight.title, payload)).toBe('');
    expect(resolved(insight.text, payload)).toBe('');
    expect(resolved(upsell.title, payload)).toBe('');
    expect(resolved(upsell.text, payload)).toBe('');
    expect(resolved(upsell.callout, payload)).toBe('');
    expect(resolved(outro.title, payload)).toBe('Это был ваш год на Авито');
    expect(resolved(outro.subtitle, payload)).toBe('Поделитесь итогами или загляните на главную');
  });

  it('uses narrative comparison as upsell callout when item callout is absent', () => {
    const payload = makePayload({
      narrative: { scenes: { offer: { comparison: 'Выгода {{value}}' } } },
      story: [
        {
          id: 'offer',
          type: ESceneType.Upsell,
          value: 'saving',
        },
      ],
    });
    const upsell = buildScenesFromStory(payload)[0];
    if (upsell.type !== ESceneType.Upsell) return;

    expect(resolved(upsell.callout, payload)).toBe(
      `Выгода ${contextFor(payload).format.currency(1250)}`,
    );
  });

  it('rejects duplicate story ids', () => {
    const payload = makePayload({
      story: [
        { id: 'same', type: ESceneType.Custom },
        { id: 'same', type: ESceneType.Custom },
      ],
    });

    expect(() => buildScenesFromStory(payload)).toThrow('Duplicate scene id: same');
  });
});
