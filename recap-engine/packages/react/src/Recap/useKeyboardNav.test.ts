import { jest } from '@jest/globals';
import { EPlayerPhase, ERecapEventType } from '@recap-engine/core';
import { act, renderHook } from '@testing-library/react';
import { createRecapValue, recapWrapper } from '../test/recapTestUtils';
import { useKeyboardNav } from './useKeyboardNav';

describe('useKeyboardNav', () => {
  it.each(['ArrowRight', ' '])('advances and prevents default for %p', (key) => {
    const next = jest.fn();
    const value = createRecapValue({
      next,
      player: { index: 0, total: 2, direction: 1, phase: EPlayerPhase.Active },
    });
    renderHook(() => useKeyboardNav(), { wrapper: recapWrapper(value) });
    const event = new KeyboardEvent('keydown', { key, cancelable: true });

    act(() => window.dispatchEvent(event));

    expect(event.defaultPrevented).toBe(true);
    expect(next).toHaveBeenCalledTimes(1);
  });

  it('goes back on ArrowLeft and ignores unrelated keys', () => {
    const prev = jest.fn();
    const next = jest.fn();
    const value = createRecapValue({ prev, next });
    renderHook(() => useKeyboardNav(), { wrapper: recapWrapper(value) });

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowLeft' }));
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter' }));
    });

    expect(prev).toHaveBeenCalledTimes(1);
    expect(next).not.toHaveBeenCalled();
  });

  it('emits complete while handling forward navigation at the end', () => {
    const onEvent = jest.fn();
    const next = jest.fn();
    const value = createRecapValue({
      onEvent,
      next,
      player: { index: 1, total: 2, direction: 1, phase: EPlayerPhase.Active },
    });
    renderHook(() => useKeyboardNav(), { wrapper: recapWrapper(value) });

    act(() => window.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight' })));

    expect(onEvent).toHaveBeenCalledWith({ type: ERecapEventType.Complete });
    expect(next).toHaveBeenCalledTimes(1);
  });

  it('removes the keydown listener on unmount', () => {
    const removeSpy = jest.spyOn(window, 'removeEventListener');
    const value = createRecapValue();
    const { unmount } = renderHook(() => useKeyboardNav(), {
      wrapper: recapWrapper(value),
    });

    unmount();

    expect(removeSpy).toHaveBeenCalledWith('keydown', expect.any(Function));
  });
});
