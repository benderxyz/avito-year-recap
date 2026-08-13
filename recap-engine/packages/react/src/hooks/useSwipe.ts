import { type PointerEvent as ReactPointerEvent, useCallback, useMemo, useRef } from 'react';

const DEFAULT_THRESHOLD = 50;
const DEFAULT_TAP_RATIO = 0.33;

export type TapNavConfig = boolean | { leftRatio?: number; rightRatio?: number };

export type SwipeOptions = {
  threshold?: number;
  tapNav?: TapNavConfig;
  holdToPause?: boolean;
  onHoldPauseChange?: (paused: boolean) => void;
};

type SwipeHandlers = {
  onPointerDown: (e: ReactPointerEvent) => void;
  onPointerUp: (e: ReactPointerEvent) => void;
  onPointerCancel: (e: ReactPointerEvent) => void;
};

function resolveTapNav(tapNav?: TapNavConfig) {
  if (!tapNav) return null;

  if (tapNav === true) {
    return { leftRatio: DEFAULT_TAP_RATIO, rightRatio: DEFAULT_TAP_RATIO };
  }

  return {
    leftRatio: tapNav.leftRatio ?? DEFAULT_TAP_RATIO,
    rightRatio: tapNav.rightRatio ?? DEFAULT_TAP_RATIO,
  };
}

function resolveOptions(thresholdOrOptions?: number | SwipeOptions): SwipeOptions {
  if (typeof thresholdOrOptions === 'number') {
    return { threshold: thresholdOrOptions };
  }

  return thresholdOrOptions ?? {};
}

function isInteractiveTarget(target: EventTarget | null): boolean {
  if (!(target instanceof Element)) return false;

  return Boolean(
    target.closest('button, a, input, textarea, select, label, [role="button"], .recap-progress'),
  );
}

export function useSwipe(
  enabled: boolean,
  onSwipeLeft: () => void,
  onSwipeRight: () => void,
  thresholdOrOptions: number | SwipeOptions = DEFAULT_THRESHOLD,
): SwipeHandlers {
  const {
    threshold = DEFAULT_THRESHOLD,
    tapNav,
    holdToPause = false,
    onHoldPauseChange,
  } = resolveOptions(thresholdOrOptions);
  const tapConfig = resolveTapNav(tapNav);
  const startX = useRef<number | null>(null);
  const isHolding = useRef(false);
  const pointerIdRef = useRef<number | null>(null);
  const targetRef = useRef<HTMLElement | null>(null);

  const releaseHold = useCallback(() => {
    if (!isHolding.current) return;

    isHolding.current = false;
    onHoldPauseChange?.(false);

    if (targetRef.current && pointerIdRef.current !== null) {
      try {
        targetRef.current.releasePointerCapture(pointerIdRef.current);
      } catch {
        // ignore release errors for detached nodes
      }
    }

    pointerIdRef.current = null;
    targetRef.current = null;
  }, [onHoldPauseChange]);

  const onPointerDown = useCallback(
    (e: ReactPointerEvent) => {
      if (!enabled || isInteractiveTarget(e.target)) return;

      startX.current = e.clientX;

      if (!holdToPause) return;

      isHolding.current = true;
      pointerIdRef.current = e.pointerId;
      targetRef.current = e.currentTarget as HTMLElement;

      try {
        targetRef.current.setPointerCapture(e.pointerId);
      } catch {
        // ignore capture errors
      }

      onHoldPauseChange?.(true);
    },
    [enabled, holdToPause, onHoldPauseChange],
  );

  const onPointerUp = useCallback(
    (e: ReactPointerEvent) => {
      if (holdToPause) {
        releaseHold();
      }

      if (!enabled || startX.current === null || isInteractiveTarget(e.target)) {
        startX.current = null;
        return;
      }

      const dx = e.clientX - startX.current;
      startX.current = null;

      if (dx <= -threshold) {
        onSwipeLeft();
        return;
      }

      if (dx >= threshold) {
        onSwipeRight();
        return;
      }

      if (!tapConfig || Math.abs(dx) >= threshold) return;

      const rect = e.currentTarget.getBoundingClientRect();
      const relativeX = e.clientX - rect.left;

      if (relativeX < rect.width * tapConfig.leftRatio) {
        onSwipeRight();
        return;
      }

      if (relativeX > rect.width * (1 - tapConfig.rightRatio)) {
        onSwipeLeft();
      }
    },
    [enabled, holdToPause, onSwipeLeft, onSwipeRight, releaseHold, tapConfig, threshold],
  );

  const onPointerCancel = useCallback(
    (_e: ReactPointerEvent) => {
      if (holdToPause) {
        releaseHold();
      }

      startX.current = null;
    },
    [holdToPause, releaseHold],
  );

  return useMemo(
    () => ({ onPointerDown, onPointerUp, onPointerCancel }),
    [onPointerCancel, onPointerDown, onPointerUp],
  );
}
