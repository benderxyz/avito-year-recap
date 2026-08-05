import { jest } from '@jest/globals';
import { act, renderHook } from '@testing-library/react';
import { usePrefersReducedMotion } from './usePrefersReducedMotion';

describe('usePrefersReducedMotion', () => {
  it.each([true, false])('returns the explicit %s override without querying media', (override) => {
    const matchMedia = jest.spyOn(window, 'matchMedia');

    const { result } = renderHook(() => usePrefersReducedMotion(override));

    expect(result.current).toBe(override);
    expect(matchMedia).not.toHaveBeenCalled();
  });

  it('reads the media query and responds to changes', () => {
    let matches = false;
    let changeListener: (() => void) | undefined;
    const addEventListener = jest.fn(
      (_type: string, listener: EventListenerOrEventListenerObject) => {
        changeListener = listener as () => void;
      },
    );
    jest.spyOn(window, 'matchMedia').mockImplementation(
      () =>
        ({
          get matches() {
            return matches;
          },
          addEventListener,
          removeEventListener: jest.fn(),
        }) as unknown as MediaQueryList,
    );

    const { result } = renderHook(() => usePrefersReducedMotion());
    expect(result.current).toBe(false);

    matches = true;
    act(() => changeListener?.());

    expect(result.current).toBe(true);
    expect(addEventListener).toHaveBeenCalledWith('change', expect.any(Function));
  });

  it('removes the change listener on unmount', () => {
    const addEventListener = jest.fn();
    const removeEventListener = jest.fn();
    jest.spyOn(window, 'matchMedia').mockReturnValue({
      matches: false,
      addEventListener,
      removeEventListener,
    } as unknown as MediaQueryList);
    const { unmount } = renderHook(() => usePrefersReducedMotion());
    const listener = addEventListener.mock.calls[0]?.[1];

    unmount();

    expect(removeEventListener).toHaveBeenCalledWith('change', listener);
  });
});
