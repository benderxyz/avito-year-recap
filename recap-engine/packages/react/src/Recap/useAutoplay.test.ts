import { afterEach, beforeEach, jest } from '@jest/globals';
import { EPlayerPhase, ERecapEventType } from '@recap-engine/core';
import { act, renderHook } from '@testing-library/react';
import { createRecapValue, recapWrapper } from '../test/recapTestUtils';
import { useAutoplay } from './useAutoplay';

describe('useAutoplay', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.clearAllTimers();
    jest.useRealTimers();
  });

  it('advances after the configured delay', () => {
    const next = jest.fn();
    const value = createRecapValue({
      next,
      player: { index: 0, total: 2, direction: 1, phase: EPlayerPhase.Active },
    });

    renderHook(() => useAutoplay({ autoplay: { delayMs: 250 } }), {
      wrapper: recapWrapper(value),
    });

    act(() => jest.advanceTimersByTime(249));
    expect(next).not.toHaveBeenCalled();

    act(() => jest.advanceTimersByTime(1));
    expect(next).toHaveBeenCalledTimes(1);
  });

  it('uses duration as the boolean-autoplay delay', () => {
    const next = jest.fn();
    const value = createRecapValue({
      next,
      player: { index: 0, total: 2, direction: 1, phase: EPlayerPhase.Active },
    });

    renderHook(() => useAutoplay({ autoplay: true, durationMs: 500 }), {
      wrapper: recapWrapper(value),
    });

    act(() => jest.advanceTimersByTime(500));
    expect(next).toHaveBeenCalledTimes(1);
  });

  it.each([
    ['disabled', { autoplay: false }, { isAnimating: false, blockMotionDone: true }],
    ['animation', { autoplay: true }, { isAnimating: true, blockMotionDone: true }],
    ['block motion', { autoplay: true }, { isAnimating: false, blockMotionDone: false }],
  ])('does not schedule while %s prevents autoplay', (_label, options, state) => {
    const value = createRecapValue(state);

    renderHook(() => useAutoplay(options), { wrapper: recapWrapper(value) });

    expect(jest.getTimerCount()).toBe(0);
  });

  it('emits complete and loops from the final scene', () => {
    const onEvent = jest.fn();
    const goTo = jest.fn();
    const next = jest.fn();
    const value = createRecapValue({
      onEvent,
      goTo,
      next,
      player: { index: 1, total: 2, direction: 1, phase: EPlayerPhase.Active },
    });

    renderHook(() => useAutoplay({ autoplay: true, loop: true, durationMs: 100 }), {
      wrapper: recapWrapper(value),
    });

    act(() => jest.advanceTimersByTime(100));

    expect(onEvent).toHaveBeenCalledWith({ type: ERecapEventType.Complete });
    expect(goTo).toHaveBeenCalledWith(0);
    expect(next).not.toHaveBeenCalled();
  });

  it('completes without looping when loop is disabled', () => {
    const onEvent = jest.fn();
    const goTo = jest.fn();
    const value = createRecapValue({
      onEvent,
      goTo,
      player: { index: 0, total: 1, direction: 1, phase: EPlayerPhase.Active },
    });

    renderHook(() => useAutoplay({ autoplay: true, durationMs: 100 }), {
      wrapper: recapWrapper(value),
    });

    act(() => jest.advanceTimersByTime(100));

    expect(onEvent).toHaveBeenCalledWith({ type: ERecapEventType.Complete });
    expect(goTo).not.toHaveBeenCalled();
  });

  it('clears its timer on unmount', () => {
    const value = createRecapValue({
      player: { index: 0, total: 2, direction: 1, phase: EPlayerPhase.Active },
    });
    const { unmount } = renderHook(() => useAutoplay({ autoplay: true }), {
      wrapper: recapWrapper(value),
    });

    expect(jest.getTimerCount()).toBe(1);
    unmount();
    expect(jest.getTimerCount()).toBe(0);
  });
});
