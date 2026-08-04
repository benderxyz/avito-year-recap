import {
  EMotionPreset,
  type MotionsConfig,
  normalizeMotion,
  type SceneDefinition,
} from '@recap-engine/core';
import { sceneMotionClass } from './sceneMotionClass';

export function useScenePresentation<TData>(
  scene: SceneDefinition<TData> | undefined,
  direction: 1 | -1,
  motions: MotionsConfig,
  reducedMotion: boolean,
) {
  const motion = normalizeMotion(
    scene?.motion ?? (direction === -1 ? motions.back : motions.default) ?? EMotionPreset.SlideUp,
    typeof motions.default === 'string' ? motions.default : EMotionPreset.SlideUp,
  );

  const motionClass = sceneMotionClass(
    motion.enter,
    direction,
    reducedMotion,
    motions.reducedMotion ?? EMotionPreset.Fade,
  );

  return { motion, motionClass };
}
