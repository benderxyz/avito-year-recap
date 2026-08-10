import { ESceneActionType, ESceneType, EShareKind, type RecapPayload } from '@recap-engine/core';
import { describe, expect, it } from 'vitest';
import { enrichShareActionsWithUrl } from '@/entities/recap';

const shareUrl = 'https://recaps.hakolr.dev/share/token';

const basePayload = {
  schemaVersion: 1 as const,
  meta: {
    vertical: 'marketplace',
    year: 2026,
    locale: 'ru-RU',
    user: { id: 'avito-42', displayName: 'Anna' },
    generatedAt: '2026-01-01T00:00:00Z',
  },
  metrics: {},
  story: [
    {
      id: 'outro',
      type: ESceneType.Outro,
      actions: [
        {
          type: ESceneActionType.Share as const,
          label: 'Поделиться',
          share: {
            kind: EShareKind.Link as const,
            title: 'Мои итоги',
            text: 'Посмотрите recap',
          },
        },
      ],
    },
  ],
} satisfies RecapPayload;

describe('enrichShareActionsWithUrl', () => {
  it('injects features.shareUrl into outro share action', () => {
    const payload = {
      ...basePayload,
      features: {
        shareEnabled: true,
        shareUrl,
      },
    } as RecapPayload;

    const enriched = enrichShareActionsWithUrl(payload);
    const outro = enriched.story[0];

    expect(outro.actions?.[0]).toMatchObject({
      type: ESceneActionType.Share,
      share: {
        url: shareUrl,
      },
    });
  });

  it('keeps existing share.url untouched', () => {
    const existingUrl = 'https://recaps.hakolr.dev/share/existing';
    const payload = {
      ...basePayload,
      features: {
        shareEnabled: true,
        shareUrl,
      },
      story: [
        {
          ...basePayload.story[0],
          actions: [
            {
              type: ESceneActionType.Share as const,
              label: 'Поделиться',
              share: {
                kind: EShareKind.Link as const,
                url: existingUrl,
              },
            },
          ],
        },
      ],
    } as RecapPayload;

    const enriched = enrichShareActionsWithUrl(payload);
    expect(enriched.story[0].actions?.[0]).toMatchObject({
      share: { url: existingUrl },
    });
  });

  it('returns payload unchanged when shareUrl is missing', () => {
    expect(enrichShareActionsWithUrl(basePayload)).toEqual(basePayload);
  });
});
