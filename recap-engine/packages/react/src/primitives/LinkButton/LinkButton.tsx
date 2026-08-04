import { EButtonVariant, ELinkTarget, ERecapEventType, ESceneActionType } from '@recap-engine/core';
import { useCallback } from 'react';
import { useRecap } from '../../context/RecapContext';
import { variantClass } from '../ActionButton/variantClass';
import type { LinkButtonProps } from './interface';

export function LinkButton({
  label,
  href,
  target = ELinkTarget.Self,
  variant = EButtonVariant.Primary,
}: LinkButtonProps) {
  const { player, onEvent, scenes } = useRecap();
  const scene = scenes[player.index];

  const handleClick = useCallback(() => {
    onEvent?.({
      type: ERecapEventType.Action,
      sceneId: scene.id,
      action: { type: ESceneActionType.Link, label, href, target, variant },
    });
  }, [href, label, onEvent, scene.id, target, variant]);

  return (
    <a
      className={variantClass(variant)}
      href={href}
      target={target}
      rel={target === ELinkTarget.Blank ? 'noreferrer' : undefined}
      onClick={handleClick}
    >
      {label}
    </a>
  );
}
