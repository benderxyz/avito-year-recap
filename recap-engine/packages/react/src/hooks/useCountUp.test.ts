import { jest } from '@jest/globals';
import { act, renderHook } from '@testing-library/react';
import { useCountUp } from './useCountUp';

describe('useCountUp', () => {
  it('returns the target and completes immediately when disabled', () => {
    const onComplete = jest.fn();
    const { result } = renderHook(() => useCountUp(42, false, 1000, onComplete));

    expect(result.current).toBe(42);
    expect(onComplete).toHaveBeenCalledTimes(1);
  });

  it('animates with easing and lands exactly on the target', () => {
    const frames: FrameRequestCallback[] = [];
    jest.spyOn(performance, 'now').mockReturnValue(100);
    jest
      .spyOn(window, 'requestAnimationFrame')
      .mockImplementation((callback: FrameRequestCallback) => {
        frames.push(callback);
        return frames.length;
      });
    const onComplete = jest.fn();
    const { result } = renderHook(() => useCountUp(100, true, 1000, onComplete));

    expect(result.current).toBe(0);

    act(() => frames.shift()?.(600));
    expect(result.current).toBe(88);
    expect(onComplete).not.toHaveBeenCalled();

    act(() => frames.shift()?.(1100));
    expect(result.current).toBe(100);
    expect(onComplete).toHaveBeenCalledTimes(1);
  });

  it('cancels the pending frame on unmount', () => {
    jest.spyOn(performance, 'now').mockReturnValue(0);
    jest.spyOn(window, 'requestAnimationFrame').mockReturnValue(17);
    const cancel = jest.spyOn(window, 'cancelAnimationFrame');
    const { unmount } = renderHook(() => useCountUp(10, true));

    unmount();

    expect(cancel).toHaveBeenCalledWith(17);
  });

  it('restarts from zero when the target changes', () => {
    jest.spyOn(performance, 'now').mockReturnValue(0);
    jest.spyOn(window, 'requestAnimationFrame').mockReturnValue(1);
    const { result, rerender } = renderHook(({ target }) => useCountUp(target, true), {
      initialProps: { target: 10 },
    });

    rerender({ target: 20 });

    expect(result.current).toBe(0);
  });
});
