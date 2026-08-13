import { jest } from '@jest/globals';
import {
  createFormatters,
  createTheme,
  EMotionPreset,
  EPlayerPhase,
  ESceneType,
  type SceneDefinition,
} from '@recap-engine/core';
import { type RenderOptions, type RenderResult, render } from '@testing-library/react';
import type { ReactElement } from 'react';
import { type RecapContextValue, RecapProvider } from '../context/RecapContext';

const defaultScenes: SceneDefinition[] = [
  { id: 'scene-1', type: ESceneType.Custom, actions: [] },
  { id: 'scene-2', type: ESceneType.Custom, actions: [] },
];

export function createRecapValue<TData = unknown>(
  overrides: Partial<RecapContextValue<TData>> = {},
): RecapContextValue<TData> {
  const scenes = overrides.scenes ?? (defaultScenes as SceneDefinition<TData>[]);
  const player = {
    index: 0,
    direction: 1 as const,
    phase: EPlayerPhase.Active,
    total: scenes.length,
    ...overrides.player,
  };

  return {
    data: {} as TData,
    theme: createTheme(),
    scenes,
    dispatch: jest.fn(),
    format: createFormatters('ru-RU'),
    motions: {
      default: EMotionPreset.SlideUp,
      back: EMotionPreset.SlideLeft,
      reducedMotion: EMotionPreset.Fade,
    },
    reducedMotion: true,
    badges: [],
    next: jest.fn(),
    prev: jest.fn(),
    goTo: jest.fn(),
    progress: 0,
    isAnimating: false,
    setAnimating: jest.fn(),
    notifyBlockMotionComplete: jest.fn(),
    blockMotionDone: true,
    isAutoplayPaused: false,
    ...overrides,
    player,
  };
}

export function renderWithRecap<TData = unknown>(
  ui: ReactElement,
  overrides: Partial<RecapContextValue<TData>> = {},
  options?: Omit<RenderOptions, 'wrapper'>,
): RenderResult & { value: RecapContextValue<TData> } {
  const value = createRecapValue(overrides);
  return {
    value,
    ...render(<RecapProvider value={value}>{ui}</RecapProvider>, options),
  };
}
