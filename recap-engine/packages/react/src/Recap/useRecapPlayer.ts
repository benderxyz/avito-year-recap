import {
  createInitialPlayerState,
  EPlayerActionType,
  ERecapEventType,
  playerReducer,
  type RecapEvent,
  type SceneDefinition,
} from '@recap-engine/core';
import { useCallback, useEffect, useReducer } from 'react';

export function useRecapPlayer<TData>(
  scenes: SceneDefinition<TData>[],
  initialSceneId: string | undefined,
  onEvent?: (event: RecapEvent) => void,
) {
  const initialIndex = Math.max(
    0,
    initialSceneId ? scenes.findIndex((s) => s.id === initialSceneId) : 0,
  );

  const [player, dispatch] = useReducer(
    playerReducer,
    createInitialPlayerState(scenes.length, initialIndex === -1 ? 0 : initialIndex),
  );

  useEffect(() => {
    dispatch({ type: EPlayerActionType.SetTotal, total: scenes.length });
  }, [scenes.length]);

  const next = useCallback(() => {
    if (player.index >= player.total - 1) {
      onEvent?.({ type: ERecapEventType.Complete });
      return;
    }
    dispatch({ type: EPlayerActionType.Next });
  }, [onEvent, player.index, player.total]);

  const prev = useCallback(() => {
    dispatch({ type: EPlayerActionType.Prev });
  }, []);

  const goTo = useCallback((index: number) => {
    dispatch({ type: EPlayerActionType.GoTo, index });
  }, []);

  const progress = player.total <= 1 ? 1 : player.index / (player.total - 1);

  return { player, dispatch, next, prev, goTo, progress };
}
