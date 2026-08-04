export enum EMotionPreset {
  Fade = 'fade',
  SlideUp = 'slide-up',
  SlideLeft = 'slide-left',
  ScaleFade = 'scale-fade',
  CountUp = 'count-up',
  BadgePop = 'badge-pop',
  StaggerText = 'stagger-text',
  CalloutIn = 'callout-in',
  None = 'none',
}

export enum EMotionEase {
  Standard = 'standard',
  Emphasize = 'emphasize',
  Exit = 'exit',
}

export type MotionConfig = {
  enter?: EMotionPreset;
  exit?: EMotionPreset;
  durationMs?: number;
  staggerMs?: number;
  ease?: EMotionEase;
};

export type MotionsConfig = {
  default?: EMotionPreset | MotionConfig;
  back?: EMotionPreset | MotionConfig;
  reducedMotion?: EMotionPreset;
};
