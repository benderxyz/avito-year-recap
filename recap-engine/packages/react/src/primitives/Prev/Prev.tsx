import { EButtonVariant, ERecapEventType, ESceneActionType } from '@recap-engine/core';
import { useCallback } from 'react';
import { useRecap } from '../../context/RecapContext';
import { ActionButton } from '../ActionButton';
import type { PrevProps } from './interface';

export function Prev({ label = 'Назад', variant = EButtonVariant.Ghost }: PrevProps) {
  const { prev, player, onEvent, scenes } = useRecap();
  const scene = scenes[player.index];

  const handleClick = useCallback(() => {
    onEvent?.({
      type: ERecapEventType.Action,
      sceneId: scene.id,
      action: { type: ESceneActionType.Prev, label, variant },
    });
    prev();
  }, [label, onEvent, prev, scene.id, variant]);

  return (
    <ActionButton variant={variant} onClick={handleClick}>
      {label}
    </ActionButton>
  );
}
