import { type PointerEvent as ReactPointerEvent, useCallback, useMemo, useRef } from 'react';

type SwipeHandlers = {
  onPointerDown: (e: ReactPointerEvent) => void;
  onPointerUp: (e: ReactPointerEvent) => void;
  onPointerCancel: () => void;
};

export function useSwipe(
  enabled: boolean,
  onSwipeLeft: () => void,
  onSwipeRight: () => void,
  threshold = 50,
): SwipeHandlers {
  const startX = useRef<number | null>(null);

  const onPointerDown = useCallback(
    (e: ReactPointerEvent) => {
      if (!enabled) return;
      startX.current = e.clientX;
    },
    [enabled],
  );

  const onPointerUp = useCallback(
    (e: ReactPointerEvent) => {
      if (!enabled || startX.current === null) return;

      const dx = e.clientX - startX.current;
      startX.current = null;

      if (dx <= -threshold) onSwipeLeft();
      else if (dx >= threshold) onSwipeRight();
    },
    [enabled, onSwipeLeft, onSwipeRight, threshold],
  );

  const onPointerCancel = useCallback(() => {
    startX.current = null;
  }, []);

  return useMemo(
    () => ({ onPointerDown, onPointerUp, onPointerCancel }),
    [onPointerCancel, onPointerDown, onPointerUp],
  );
}
