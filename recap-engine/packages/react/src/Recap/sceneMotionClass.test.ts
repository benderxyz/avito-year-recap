import { describe, expect, it } from '@jest/globals';
import { EMotionPreset } from '@recap-engine/core';
import { sceneMotionClass } from './sceneMotionClass';

describe('sceneMotionClass', () => {
  it('reverses horizontal slide direction when navigating back', () => {
    expect(sceneMotionClass(EMotionPreset.SlideLeft, 1, false, EMotionPreset.Fade)).toBe(
      'recap-motion--slide-left',
    );
    expect(sceneMotionClass(EMotionPreset.SlideLeft, -1, false, EMotionPreset.Fade)).toBe(
      'recap-motion--slide-right',
    );
  });

  it.each([
    EMotionPreset.CountUp,
    EMotionPreset.BadgePop,
    EMotionPreset.StaggerText,
    EMotionPreset.CalloutIn,
  ])('maps block preset %s to fade', (preset: EMotionPreset) => {
    expect(sceneMotionClass(preset, 1, false, EMotionPreset.None)).toBe('recap-motion--fade');
  });

  it('uses the reduced-motion fallback and supports none', () => {
    expect(sceneMotionClass(EMotionPreset.SlideUp, 1, true, EMotionPreset.None)).toBe(
      'recap-motion--none',
    );
  });

  it('builds the class for other presets', () => {
    expect(sceneMotionClass(EMotionPreset.ScaleFade, 1, false, EMotionPreset.Fade)).toBe(
      'recap-motion--scale-fade',
    );
  });
});
