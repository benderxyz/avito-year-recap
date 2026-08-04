import { ESceneType, type SceneDefinition } from '@recap-engine/core';
import { useCallback, useEffect, useState } from 'react';

const BLOCK_MOTION_SCENE_TYPES = new Set<ESceneType>([
  ESceneType.Stat,
  ESceneType.Achievement,
  ESceneType.Insight,
  ESceneType.Upsell,
  ESceneType.Blocks,
]);

const BLOCK_MOTION_FALLBACK_MS = 1400;

function needsBlockMotionWait<TData>(
  scene: SceneDefinition<TData> | undefined,
  reducedMotion: boolean,
): boolean {
  if (!scene || reducedMotion) return false;
  return BLOCK_MOTION_SCENE_TYPES.has(scene.type);
}

export function useBlockMotion<TData>(
  scene: SceneDefinition<TData> | undefined,
  reducedMotion: boolean,
) {
  const [blockMotionDone, setBlockMotionDone] = useState(true);

  useEffect(() => {
    setBlockMotionDone(false);

    if (!needsBlockMotionWait(scene, reducedMotion)) {
      setBlockMotionDone(true);
      return;
    }

    const fallback = window.setTimeout(() => setBlockMotionDone(true), BLOCK_MOTION_FALLBACK_MS);
    return () => window.clearTimeout(fallback);
  }, [scene, reducedMotion]);

  const notifyBlockMotionComplete = useCallback(() => {
    setBlockMotionDone(true);
  }, []);

  return { blockMotionDone, notifyBlockMotionComplete };
}
