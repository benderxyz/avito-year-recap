import { describe, expect, it } from '@jest/globals';
import { EMotionPreset } from '@recap-engine/core';
import { renderHook } from '@testing-library/react';
import { introScene } from '../test/recapTestUtils';
import { useScenePresentation } from './useScenePresentation';

describe('useScenePresentation', () => {
  it('prefers scene motion and normalizes its configuration', () => {
    const scene = {
      ...introScene('scene'),
      motion: { enter: EMotionPreset.ScaleFade, durationMs: 700 },
    };
    const { result } = renderHook(() =>
      useScenePresentation(scene, 1, { default: EMotionPreset.SlideUp }, false),
    );

    expect(result.current.motion).toMatchObject({
      enter: EMotionPreset.ScaleFade,
      durationMs: 700,
    });
    expect(result.current.motionClass).toBe('recap-motion--scale-fade');
  });

  it('uses back motion when navigating backward', () => {
    const { result } = renderHook(() =>
      useScenePresentation(
        introScene('scene'),
        -1,
        { default: EMotionPreset.Fade, back: EMotionPreset.SlideLeft },
        false,
      ),
    );

    expect(result.current.motion.enter).toBe(EMotionPreset.SlideLeft);
    expect(result.current.motionClass).toBe('recap-motion--slide-right');
  });

  it('falls back to slide-up and applies reduced motion', () => {
    const { result } = renderHook(() => useScenePresentation(undefined, 1, {}, true));

    expect(result.current.motion.enter).toBe(EMotionPreset.SlideUp);
    expect(result.current.motionClass).toBe('recap-motion--fade');
  });

  it('uses a configured default as the fallback for partial motion objects', () => {
    const scene = {
      ...introScene('scene'),
      motion: { durationMs: 250 },
    };
    const { result } = renderHook(() =>
      useScenePresentation(scene, 1, { default: EMotionPreset.ScaleFade }, false),
    );

    expect(result.current.motion.enter).toBe(EMotionPreset.ScaleFade);
  });
});
