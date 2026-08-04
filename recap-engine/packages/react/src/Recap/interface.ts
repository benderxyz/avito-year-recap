import type {
  Badge,
  MotionsConfig,
  RecapEvent,
  ResolvedTheme,
  SceneDefinition,
} from '@recap-engine/core';
import type { ReactNode } from 'react';

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
  reducedMotion?: boolean;
  className?: string;
  slots?: {
    background?: ReactNode;
    header?: ReactNode;
    footer?: ReactNode;
  };
  onEvent?: (event: RecapEvent) => void;
};
