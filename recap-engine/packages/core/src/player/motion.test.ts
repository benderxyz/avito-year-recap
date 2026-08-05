import { describe, expect, it } from '@jest/globals';

import { EMotionPreset } from '../types/motion';
import { DEFAULT_MOTION_DURATION_MS, normalizeMotion } from './motion';

describe('normalizeMotion', () => {
  it('uses the default fallback when motion is absent', () => {
    expect(normalizeMotion(undefined)).toEqual({
      enter: EMotionPreset.SlideUp,
      exit: EMotionPreset.SlideUp,
      durationMs: DEFAULT_MOTION_DURATION_MS,
    });
  });

  it('uses a custom fallback when motion is absent', () => {
    expect(normalizeMotion(undefined, EMotionPreset.Fade)).toEqual({
      enter: EMotionPreset.Fade,
      exit: EMotionPreset.Fade,
      durationMs: DEFAULT_MOTION_DURATION_MS,
    });
  });

  it('expands a preset to matching enter and exit motion', () => {
    expect(normalizeMotion(EMotionPreset.ScaleFade)).toEqual({
      enter: EMotionPreset.ScaleFade,
      exit: EMotionPreset.ScaleFade,
      durationMs: DEFAULT_MOTION_DURATION_MS,
    });
  });

  it('normalizes a complete motion config', () => {
    expect(
      normalizeMotion({
        enter: EMotionPreset.BadgePop,
        exit: EMotionPreset.Fade,
        durationMs: 900,
      }),
    ).toEqual({
      enter: EMotionPreset.BadgePop,
      exit: EMotionPreset.Fade,
      durationMs: 900,
    });
  });

  it('fills omitted config properties without replacing zero duration', () => {
    expect(
      normalizeMotion({ exit: EMotionPreset.None, durationMs: 0 }, EMotionPreset.SlideLeft),
    ).toEqual({
      enter: EMotionPreset.SlideLeft,
      exit: EMotionPreset.None,
      durationMs: 0,
    });
  });
});
