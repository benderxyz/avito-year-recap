import { jest } from '@jest/globals';
import { EPlayerPhase, ESceneType, type SceneDefinition } from '@recap-engine/core';
import type { PropsWithChildren } from 'react';
import { type RecapContextValue, RecapProvider } from '../context/RecapProvider/RecapProvider';

export const introScene = (id: string): SceneDefinition<unknown> => ({
  id,
  type: ESceneType.Intro,
  title: id,
});

export const statScene = (id: string): SceneDefinition<unknown> => ({
  id,
  type: ESceneType.Stat,
  value: 1,
});

export function createRecapValue(overrides: Partial<RecapContextValue> = {}): RecapContextValue {
  return {
    data: {},
    theme: {} as RecapContextValue['theme'],
    scenes: [],
    player: {
      index: 0,
      direction: 1,
      phase: EPlayerPhase.Active,
      total: 1,
    },
    dispatch: jest.fn(),
    format: {} as RecapContextValue['format'],
    motions: {},
    reducedMotion: false,
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
  };
}

export function recapWrapper(value: RecapContextValue) {
  return function RecapTestWrapper({ children }: PropsWithChildren) {
    return <RecapProvider value={value}>{children}</RecapProvider>;
  };
}
