import {
  EPlayerActionType,
  EPlayerPhase,
  ERecapEventType,
  type PlayerAction,
  type RecapEvent,
  type SceneDefinition,
} from '@recap-engine/core';
import { type Dispatch, useEffect, useRef } from 'react';

type Options<TData> = {
  scene: SceneDefinition<TData> | undefined;
  scenes: SceneDefinition<TData>[];
  index: number;
  durationMs: number;
  reducedMotion: boolean;
  onEvent?: (event: RecapEvent) => void;
  setAnimating: (value: boolean) => void;
  dispatch: Dispatch<PlayerAction>;
};

export function useSceneLifecycle<TData>({
  scene,
  scenes,
  index,
  durationMs,
  reducedMotion,
  onEvent,
  setAnimating,
  dispatch,
}: Options<TData>) {
  const prevIndexRef = useRef(index);

  useEffect(() => {
    const prevIndex = prevIndexRef.current;
    const prevScene = scenes[prevIndex];

    if (prevIndex !== index && prevScene) {
      onEvent?.({
        type: ERecapEventType.SceneExit,
        sceneId: prevScene.id,
        index: prevIndex,
      });
    }

    if (scene) {
      onEvent?.({
        type: ERecapEventType.SceneEnter,
        sceneId: scene.id,
        index,
      });
    }

    prevIndexRef.current = index;

    setAnimating(true);
    dispatch({ type: EPlayerActionType.SetPhase, phase: EPlayerPhase.Enter });

    const timer = window.setTimeout(
      () => {
        setAnimating(false);
        dispatch({
          type: EPlayerActionType.SetPhase,
          phase: EPlayerPhase.Active,
        });
      },
      reducedMotion ? 0 : durationMs,
    );

    return () => window.clearTimeout(timer);
  }, [dispatch, durationMs, index, onEvent, reducedMotion, scene, scenes, setAnimating]);
}
