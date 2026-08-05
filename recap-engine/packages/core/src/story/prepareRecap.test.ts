import { describe, expect, it } from '@jest/globals';

import { EMetricType, type RecapPayload } from '../types/payload';
import { ESceneType } from '../types/scenes';
import { prepareRecap } from './prepareRecap';

const payload = (overrides: Partial<RecapPayload> = {}): RecapPayload => ({
  schemaVersion: 1,
  meta: {
    vertical: 'auto',
    year: 2025,
    locale: 'ru-RU',
    user: { id: '42', displayName: 'Максим' },
    generatedAt: '2026-01-01T00:00:00Z',
  },
  metrics: { views: { type: EMetricType.Number, value: 10 } },
  story: [{ id: 'custom', type: ESceneType.Custom, sceneType: 'chart' }],
  ...overrides,
});

describe('prepareRecap', () => {
  it('keeps payload identity, builds scenes and exposes locale', () => {
    const source = payload();
    const prepared = prepareRecap(source);

    expect(prepared.data).toBe(source);
    expect(prepared.locale).toBe('ru-RU');
    expect(prepared.scenes).toEqual([
      {
        id: 'custom',
        type: ESceneType.Custom,
        sceneType: 'chart',
        props: undefined,
        motion: undefined,
        durationMs: undefined,
        actions: undefined,
      },
    ]);
  });

  it('maps backend badges to public badges without mutating the payload', () => {
    const source = payload({
      badges: [
        {
          id: 'top',
          title: 'В топе',
          description: 'Вы среди лучших',
          iconUrl: '/top.svg',
        },
        {
          id: 'active',
          title: 'Активный',
          description: 'Много действий',
        },
      ],
    });

    expect(prepareRecap(source).badges).toEqual([
      {
        id: 'top',
        title: 'В топе',
        description: 'Вы среди лучших',
        icon: '/top.svg',
      },
      {
        id: 'active',
        title: 'Активный',
        description: 'Много действий',
        icon: undefined,
      },
    ]);
    expect(source.badges?.[0].iconUrl).toBe('/top.svg');
  });

  it('defaults missing badges to an empty list', () => {
    expect(prepareRecap(payload({ badges: undefined })).badges).toEqual([]);
  });

  it('propagates duplicate scene validation', () => {
    const source = payload({
      story: [
        { id: 'duplicate', type: ESceneType.Custom },
        { id: 'duplicate', type: ESceneType.Custom },
      ],
    });

    expect(() => prepareRecap(source)).toThrow('Duplicate scene id: duplicate');
  });
});
