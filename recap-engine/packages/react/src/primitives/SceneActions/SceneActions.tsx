import {
  EButtonVariant,
  ERecapEventType,
  ESceneActionType,
  type SceneAction,
} from '@recap-engine/core';
import { useCallback } from 'react';
import { useRecap } from '../../context/RecapContext';
import { ActionButton } from '../ActionButton';
import { Actions } from '../Actions';
import { LinkButton } from '../LinkButton';
import { Next } from '../Next';
import { Prev } from '../Prev';
import { ShareButton } from '../ShareButton';
import type { SceneActionsProps } from './interface';

type GoToAction = Extract<SceneAction, { type: ESceneActionType.GoTo }>;
type CustomAction = Extract<SceneAction, { type: ESceneActionType.Custom }>;

function sceneActionKey(action: SceneAction): string {
  switch (action.type) {
    case ESceneActionType.Next:
      return `next:${action.label ?? ''}:${action.variant ?? ''}`;
    case ESceneActionType.Prev:
      return `prev:${action.label ?? ''}:${action.variant ?? ''}`;
    case ESceneActionType.Link:
      return `link:${action.href}`;
    case ESceneActionType.Share:
      return `share:${action.label}`;
    case ESceneActionType.GoTo:
      return `goto:${action.sceneId}`;
    case ESceneActionType.Custom:
      return `custom:${action.id}`;
    default: {
      const _exhaustive: never = action;
      return _exhaustive;
    }
  }
}

function GoToActionButton({ action }: { action: GoToAction }) {
  const { onEvent, scenes, player, goTo } = useRecap();
  const scene = scenes[player.index];
  const targetIndex = scenes.findIndex((item) => item.id === action.sceneId);

  const handleClick = useCallback(() => {
    onEvent?.({
      type: ERecapEventType.Action,
      sceneId: scene.id,
      action,
    });
    if (targetIndex >= 0) goTo(targetIndex);
  }, [action, goTo, onEvent, scene.id, targetIndex]);

  return (
    <ActionButton variant={action.variant ?? EButtonVariant.Ghost} onClick={handleClick}>
      {action.label}
    </ActionButton>
  );
}

function CustomActionButton({ action }: { action: CustomAction }) {
  const { onEvent, scenes, player, next, prev, goTo } = useRecap();
  const scene = scenes[player.index];

  const handleClick = useCallback(() => {
    onEvent?.({
      type: ERecapEventType.Action,
      sceneId: scene.id,
      action,
    });
    if (action.id === 'next') next();
    if (action.id === 'prev') prev();
    if (action.id === 'restart') goTo(0);
  }, [action, goTo, next, onEvent, prev, scene.id]);

  return (
    <ActionButton variant={action.variant} onClick={handleClick}>
      {action.label}
    </ActionButton>
  );
}

export function SceneActions({ actions }: SceneActionsProps) {
  if (actions === undefined) {
    return (
      <Actions>
        <Prev />
        <Next />
      </Actions>
    );
  }

  if (actions.length === 0) {
    return null;
  }

  return (
    <Actions>
      {actions.map((action) => {
        const key = sceneActionKey(action);
        switch (action.type) {
          case ESceneActionType.Next:
            return <Next key={key} label={action.label} variant={action.variant} />;
          case ESceneActionType.Prev:
            return <Prev key={key} label={action.label} variant={action.variant} />;
          case ESceneActionType.Link:
            return (
              <LinkButton
                key={key}
                label={action.label}
                href={action.href}
                target={action.target}
                variant={action.variant}
              />
            );
          case ESceneActionType.Share:
            return (
              <ShareButton
                key={key}
                label={action.label}
                share={action.share}
                variant={action.variant}
              />
            );
          case ESceneActionType.GoTo:
            return <GoToActionButton key={key} action={action} />;
          case ESceneActionType.Custom:
            return <CustomActionButton key={key} action={action} />;
          default: {
            const _exhaustive: never = action;
            return _exhaustive;
          }
        }
      })}
    </Actions>
  );
}
