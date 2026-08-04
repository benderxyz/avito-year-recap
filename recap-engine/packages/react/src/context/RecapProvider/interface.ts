import type {
  Badge,
  Formatters,
  MotionsConfig,
  PlayerAction,
  PlayerState,
  RecapEvent,
  ResolvedTheme,
  SceneDefinition,
} from '@recap-engine/core';
import type { Dispatch, PropsWithChildren } from 'react';

export type RecapContextValue<TData = unknown> = {
  data: TData;
  theme: ResolvedTheme;
  scenes: SceneDefinition<TData>[];
  player: PlayerState;
  dispatch: Dispatch<PlayerAction>;
  format: Formatters;
  motions: MotionsConfig;
  reducedMotion: boolean;
  badges: Badge[];
  onEvent?: (event: RecapEvent) => void;
  next: () => void;
  prev: () => void;
  goTo: (index: number) => void;
  progress: number;
  isAnimating: boolean;
  setAnimating: (value: boolean) => void;
  notifyBlockMotionComplete: () => void;
  blockMotionDone: boolean;
};

export type RecapProviderProps<TData = unknown> = PropsWithChildren<{
  value: RecapContextValue<TData>;
}>;
