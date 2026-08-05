import { afterEach, beforeEach, jest } from '@jest/globals';
import { EPlayerActionType, EPlayerPhase, ERecapEventType } from '@recap-engine/core';
import { act, renderHook } from '@testing-library/react';
import { introScene } from '../test/recapTestUtils';
import { useSceneLifecycle } from './useSceneLifecycle';

describe('useSceneLifecycle', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.clearAllTimers();
    jest.useRealTimers();
  });

  it('enters a scene and activates it after the duration', () => {
    const scenes = [introScene('first')];
    const onEvent = jest.fn();
    const setAnimating = jest.fn();
    const dispatch = jest.fn();

    renderHook(() =>
      useSceneLifecycle({
        scene: scenes[0],
        scenes,
        index: 0,
        durationMs: 300,
        reducedMotion: false,
        onEvent,
        setAnimating,
        dispatch,
      }),
    );

    expect(onEvent).toHaveBeenCalledWith({
      type: ERecapEventType.SceneEnter,
      sceneId: 'first',
      index: 0,
    });
    expect(setAnimating).toHaveBeenCalledWith(true);
    expect(dispatch).toHaveBeenCalledWith({
      type: EPlayerActionType.SetPhase,
      phase: EPlayerPhase.Enter,
    });

    act(() => jest.advanceTimersByTime(300));

    expect(setAnimating).toHaveBeenLastCalledWith(false);
    expect(dispatch).toHaveBeenLastCalledWith({
      type: EPlayerActionType.SetPhase,
      phase: EPlayerPhase.Active,
    });
  });

  it('emits exit before enter when the index changes', () => {
    const scenes = [introScene('first'), introScene('second')];
    const onEvent = jest.fn();
    const stable = {
      scenes,
      durationMs: 100,
      reducedMotion: false,
      onEvent,
      setAnimating: jest.fn(),
      dispatch: jest.fn(),
    };
    const { rerender } = renderHook(
      ({ index }) =>
        useSceneLifecycle({
          ...stable,
          scene: scenes[index],
          index,
        }),
      { initialProps: { index: 0 } },
    );
    onEvent.mockClear();

    rerender({ index: 1 });

    expect(onEvent.mock.calls).toEqual([
      [{ type: ERecapEventType.SceneExit, sceneId: 'first', index: 0 }],
      [{ type: ERecapEventType.SceneEnter, sceneId: 'second', index: 1 }],
    ]);
    expect(jest.getTimerCount()).toBe(1);
  });

  it('activates immediately for reduced motion', () => {
    const setAnimating = jest.fn();
    const dispatch = jest.fn();
    const scene = introScene('reduced');

    renderHook(() =>
      useSceneLifecycle({
        scene,
        scenes: [scene],
        index: 0,
        durationMs: 500,
        reducedMotion: true,
        setAnimating,
        dispatch,
      }),
    );

    act(() => jest.advanceTimersByTime(0));

    expect(setAnimating).toHaveBeenLastCalledWith(false);
    expect(dispatch).toHaveBeenLastCalledWith({
      type: EPlayerActionType.SetPhase,
      phase: EPlayerPhase.Active,
    });
  });

  it('clears the activation timer on unmount', () => {
    const scene = introScene('first');
    const { unmount } = renderHook(() =>
      useSceneLifecycle({
        scene,
        scenes: [scene],
        index: 0,
        durationMs: 100,
        reducedMotion: false,
        setAnimating: jest.fn(),
        dispatch: jest.fn(),
      }),
    );

    expect(jest.getTimerCount()).toBe(1);
    unmount();
    expect(jest.getTimerCount()).toBe(0);
  });
});
