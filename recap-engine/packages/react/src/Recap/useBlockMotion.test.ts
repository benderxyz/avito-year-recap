import { afterEach, beforeEach, expect, jest, test } from '@jest/globals';
import { act, renderHook } from '@testing-library/react';
import { introScene, statScene } from '../test/recapTestUtils';
import { useBlockMotion } from './useBlockMotion';

describe('useBlockMotion', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.clearAllTimers();
    jest.useRealTimers();
  });

  test('waits for block motion and finishes after the fallback', () => {
    const scene = statScene('stat');
    const { result } = renderHook(() => useBlockMotion(scene, false));

    expect(result.current.blockMotionDone).toBe(false);

    act(() => jest.advanceTimersByTime(1399));
    expect(result.current.blockMotionDone).toBe(false);

    act(() => jest.advanceTimersByTime(1));
    expect(result.current.blockMotionDone).toBe(true);
  });

  test('can be completed explicitly before the fallback', () => {
    const scene = statScene('stat');
    const { result } = renderHook(() => useBlockMotion(scene, false));

    act(() => result.current.notifyBlockMotionComplete());

    expect(result.current.blockMotionDone).toBe(true);
  });

  test('skips waiting for non-block scenes and reduced motion', () => {
    const { result, rerender } = renderHook(
      ({ reduced }) => useBlockMotion(introScene('intro'), reduced),
      { initialProps: { reduced: false } },
    );

    expect(result.current.blockMotionDone).toBe(true);

    rerender({ reduced: true });
    expect(result.current.blockMotionDone).toBe(true);
    expect(jest.getTimerCount()).toBe(0);
  });

  test('clears the fallback on scene change and unmount', () => {
    const { rerender, unmount } = renderHook(({ id }) => useBlockMotion(statScene(id), false), {
      initialProps: { id: 'first' },
    });

    expect(jest.getTimerCount()).toBe(1);
    rerender({ id: 'second' });
    expect(jest.getTimerCount()).toBe(1);

    unmount();
    expect(jest.getTimerCount()).toBe(0);
  });
});
