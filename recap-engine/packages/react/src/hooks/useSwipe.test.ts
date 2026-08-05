import { jest } from '@jest/globals';
import { act, renderHook } from '@testing-library/react';
import type { PointerEvent as ReactPointerEvent } from 'react';
import { useSwipe } from './useSwipe';

const pointerAt = (clientX: number) => ({ clientX }) as ReactPointerEvent;

describe('useSwipe', () => {
  it('detects left and right swipes at the threshold boundary', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const { result } = renderHook(() => useSwipe(true, onLeft, onRight, 50));

    act(() => {
      result.current.onPointerDown(pointerAt(100));
      result.current.onPointerUp(pointerAt(50));
      result.current.onPointerDown(pointerAt(100));
      result.current.onPointerUp(pointerAt(150));
    });

    expect(onLeft).toHaveBeenCalledTimes(1);
    expect(onRight).toHaveBeenCalledTimes(1);
  });

  it('ignores movement below the threshold', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const { result } = renderHook(() => useSwipe(true, onLeft, onRight, 50));

    act(() => {
      result.current.onPointerDown(pointerAt(100));
      result.current.onPointerUp(pointerAt(51));
    });

    expect(onLeft).not.toHaveBeenCalled();
    expect(onRight).not.toHaveBeenCalled();
  });

  it('does nothing while disabled or without a pointer start', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const { result } = renderHook(() => useSwipe(false, onLeft, onRight));

    act(() => {
      result.current.onPointerDown(pointerAt(100));
      result.current.onPointerUp(pointerAt(0));
    });

    expect(onLeft).not.toHaveBeenCalled();
    expect(onRight).not.toHaveBeenCalled();
  });

  it('clears the gesture on pointer cancel and after pointer up', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const { result } = renderHook(() => useSwipe(true, onLeft, onRight));

    act(() => {
      result.current.onPointerDown(pointerAt(100));
      result.current.onPointerCancel();
      result.current.onPointerUp(pointerAt(0));
      result.current.onPointerDown(pointerAt(100));
      result.current.onPointerUp(pointerAt(0));
      result.current.onPointerUp(pointerAt(0));
    });

    expect(onLeft).toHaveBeenCalledTimes(1);
    expect(onRight).not.toHaveBeenCalled();
  });
});
