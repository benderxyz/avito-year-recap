import { EMotionPreset, type MotionConfig } from '../types/motion';

type NormalizedMotion = {
  enter: EMotionPreset;
  exit: EMotionPreset;
  durationMs: number;
};

export function normalizeMotion(
  motion: EMotionPreset | MotionConfig | undefined,
  fallback: EMotionPreset = EMotionPreset.SlideUp,
): NormalizedMotion {
  if (!motion) {
    return { enter: fallback, exit: fallback, durationMs: 420 };
  }

  if (typeof motion === 'string') {
    return { enter: motion, exit: motion, durationMs: 420 };
  }

  return {
    enter: motion.enter ?? fallback,
    exit: motion.exit ?? fallback,
    durationMs: motion.durationMs ?? 420,
  };
}
