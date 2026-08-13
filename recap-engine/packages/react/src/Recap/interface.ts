import type {
  Badge,
  MotionsConfig,
  RecapEvent,
  ResolvedTheme,
  SceneDefinition,
} from '@recap-engine/core';
import type { ReactNode } from 'react';

export type TapNavConfig = boolean | { leftRatio?: number; rightRatio?: number };

export type RecapProps<TData = unknown> = {
  theme: ResolvedTheme;
  data: TData;
  scenes: SceneDefinition<TData>[];
  badges?: Badge[];
  motions?: MotionsConfig;
  locale?: string;
  initialSceneId?: string;
  autoplay?: boolean | { delayMs: number };
  loop?: boolean;
  gestures?: boolean;
  /** Tap left/right screen zones to navigate (stories-style). Default: false */
  tapNav?: TapNavConfig;
  /** Pause autoplay while pointer is held. Default: false */
  holdToPause?: boolean;
  reducedMotion?: boolean;
  className?: string;
  slots?: {
    background?: ReactNode;
    header?: ReactNode;
    footer?: ReactNode;
  };
  onEvent?: (event: RecapEvent) => void;
};
