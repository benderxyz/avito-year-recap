import { EMotionPreset, type MotionConfig } from '../types/motion';

/** Default enter/exit duration when motion config omits `durationMs`. */
export const DEFAULT_MOTION_DURATION_MS = 420;

export type NormalizedMotion = {
  enter: EMotionPreset;
  exit: EMotionPreset;
  durationMs: number;
};

export function normalizeMotion(
  motion: EMotionPreset | MotionConfig | undefined,
  fallback: EMotionPreset = EMotionPreset.SlideUp,
): NormalizedMotion {
  if (!motion) {
    return { enter: fallback, exit: fallback, durationMs: DEFAULT_MOTION_DURATION_MS };
  }

  if (typeof motion === 'string') {
    return { enter: motion, exit: motion, durationMs: DEFAULT_MOTION_DURATION_MS };
  }

  return {
    enter: motion.enter ?? fallback,
    exit: motion.exit ?? fallback,
    durationMs: motion.durationMs ?? DEFAULT_MOTION_DURATION_MS,
  };
}
