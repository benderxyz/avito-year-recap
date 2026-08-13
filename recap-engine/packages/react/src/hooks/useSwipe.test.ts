import { jest } from '@jest/globals';
import { act, renderHook } from '@testing-library/react';
import type { PointerEvent as ReactPointerEvent } from 'react';
import { useSwipe } from './useSwipe';

const pointerAt = (clientX: number, target?: EventTarget | null) =>
  ({
    clientX,
    target: target ?? document.body,
    currentTarget: document.createElement('div'),
    pointerId: 1,
  }) as unknown as ReactPointerEvent;

describe('useSwipe', () => {
  beforeEach(() => {
    document.body.innerHTML = '';
  });

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

  it('ignores movement below the threshold when tap nav is disabled', () => {
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

  it('navigates by tap zones when tapNav is enabled', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const root = document.createElement('div');
    Object.defineProperty(root, 'getBoundingClientRect', {
      value: () => ({ left: 0, width: 300, top: 0, height: 600 }),
    });

    const { result } = renderHook(() => useSwipe(true, onLeft, onRight, { tapNav: true }));

    act(() => {
      result.current.onPointerDown({ ...pointerAt(40), currentTarget: root } as ReactPointerEvent);
      result.current.onPointerUp({ ...pointerAt(42), currentTarget: root } as ReactPointerEvent);
      result.current.onPointerDown({ ...pointerAt(260), currentTarget: root } as ReactPointerEvent);
      result.current.onPointerUp({ ...pointerAt(258), currentTarget: root } as ReactPointerEvent);
    });

    expect(onRight).toHaveBeenCalledTimes(1);
    expect(onLeft).toHaveBeenCalledTimes(1);
  });

  it('does not tap-navigate in the center zone', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const root = document.createElement('div');
    Object.defineProperty(root, 'getBoundingClientRect', {
      value: () => ({ left: 0, width: 300, top: 0, height: 600 }),
    });

    const { result } = renderHook(() => useSwipe(true, onLeft, onRight, { tapNav: true }));

    act(() => {
      result.current.onPointerDown({ ...pointerAt(150), currentTarget: root } as ReactPointerEvent);
      result.current.onPointerUp({ ...pointerAt(150), currentTarget: root } as ReactPointerEvent);
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
      result.current.onPointerCancel(pointerAt(100));
      result.current.onPointerUp(pointerAt(0));
      result.current.onPointerDown(pointerAt(100));
      result.current.onPointerUp(pointerAt(0));
      result.current.onPointerUp(pointerAt(0));
    });

    expect(onLeft).toHaveBeenCalledTimes(1);
    expect(onRight).not.toHaveBeenCalled();
  });

  it('pauses and resumes autoplay while the pointer is held', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const onHoldPauseChange = jest.fn();
    const root = document.createElement('div');
    root.setPointerCapture = jest.fn();
    root.releasePointerCapture = jest.fn();

    const { result } = renderHook(() =>
      useSwipe(true, onLeft, onRight, { holdToPause: true, onHoldPauseChange }),
    );

    act(() => {
      result.current.onPointerDown({ ...pointerAt(150), currentTarget: root } as ReactPointerEvent);
    });
    expect(onHoldPauseChange).toHaveBeenCalledWith(true);

    act(() => {
      result.current.onPointerUp({ ...pointerAt(150), currentTarget: root } as ReactPointerEvent);
    });
    expect(onHoldPauseChange).toHaveBeenCalledWith(false);
  });

  it('ignores interactive targets', () => {
    const onLeft = jest.fn();
    const onRight = jest.fn();
    const onHoldPauseChange = jest.fn();
    const button = document.createElement('button');

    const { result } = renderHook(() =>
      useSwipe(true, onLeft, onRight, { tapNav: true, holdToPause: true, onHoldPauseChange }),
    );

    act(() => {
      result.current.onPointerDown(pointerAt(10, button));
      result.current.onPointerUp(pointerAt(10, button));
    });

    expect(onLeft).not.toHaveBeenCalled();
    expect(onRight).not.toHaveBeenCalled();
    expect(onHoldPauseChange).not.toHaveBeenCalled();
  });
});
