import { EButtonVariant, ERecapEventType, ESceneActionType } from '@recap-engine/core';
import { useCallback } from 'react';
import { useRecap } from '../../context/RecapContext';
import { ActionButton } from '../ActionButton';
import type { ShareButtonProps } from './interface';

export function ShareButton({
  label,
  share,
  variant = EButtonVariant.Secondary,
}: ShareButtonProps) {
  const { player, onEvent, scenes } = useRecap();
  const scene = scenes[player.index];

  const handleClick = useCallback(async () => {
    onEvent?.({
      type: ERecapEventType.Share,
      sceneId: scene.id,
      payload: share,
    });

    onEvent?.({
      type: ERecapEventType.Action,
      sceneId: scene.id,
      action: { type: ESceneActionType.Share, label, share, variant },
    });

    const payload = {
      title: share.title ?? 'Мой год',
      text: share.text ?? 'Смотрите мой recap',
      url: share.url ?? (typeof window !== 'undefined' ? window.location.href : undefined),
    };

    try {
      if (navigator.share) {
        await navigator.share(payload);
      } else if (payload.url && navigator.clipboard) {
        await navigator.clipboard.writeText(payload.url);
      }
    } catch {}
  }, [label, onEvent, scene.id, share, variant]);

  return (
    <ActionButton variant={variant} onClick={handleClick}>
      {label}
    </ActionButton>
  );
}
