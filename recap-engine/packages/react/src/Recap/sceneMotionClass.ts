import { EMotionPreset } from '@recap-engine/core';

export function sceneMotionClass(
  enter: string,
  direction: 1 | -1,
  reducedMotion: boolean,
  fallbackReduced: string,
): string {
  const preset = reducedMotion ? fallbackReduced : enter;
  if (preset === EMotionPreset.None) return 'recap-motion--none';
  if (preset === EMotionPreset.SlideLeft) {
    return direction === 1 ? 'recap-motion--slide-left' : 'recap-motion--slide-right';
  }
  if (
    preset === EMotionPreset.CountUp ||
    preset === EMotionPreset.BadgePop ||
    preset === EMotionPreset.StaggerText ||
    preset === EMotionPreset.CalloutIn
  ) {
    return 'recap-motion--fade';
  }
  return `recap-motion--${preset}`;
}
