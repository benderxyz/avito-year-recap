import { ERecapEventType } from '@recap-engine/core';
import { useCallback, useEffect, useRef } from 'react';
import { useRecap } from '../context/RecapContext';

export function useAutoplay(options: {
  autoplay?: boolean | { delayMs: number };
  loop?: boolean;
  durationMs?: number;
  isPaused?: boolean;
}) {
  const { autoplay, loop = false, durationMs, isPaused = false } = options;
  const { player, next, goTo, onEvent, isAnimating, blockMotionDone } = useRecap();
  const timerRef = useRef<number | null>(null);
  const deadlineRef = useRef<number | null>(null);
  const remainingRef = useRef<number | null>(null);

  const autoplayEnabled = Boolean(autoplay);
  const autoplayDelay = typeof autoplay === 'object' ? autoplay.delayMs : (durationMs ?? 3200);

  const clearTimer = useCallback(() => {
    if (timerRef.current !== null) {
      window.clearTimeout(timerRef.current);
      timerRef.current = null;
    }
    deadlineRef.current = null;
  }, []);

  const schedule = useCallback(
    (delayMs: number) => {
      clearTimer();
      deadlineRef.current = Date.now() + delayMs;
      timerRef.current = window.setTimeout(() => {
        timerRef.current = null;
        deadlineRef.current = null;
        remainingRef.current = null;

        if (player.index >= player.total - 1) {
          onEvent?.({ type: ERecapEventType.Complete });
          if (loop) goTo(0);
          return;
        }

        next();
      }, delayMs);
    },
    [clearTimer, goTo, loop, next, onEvent, player.index, player.total],
  );

  const sceneIndexRef = useRef(player.index);

  useEffect(() => {
    if (sceneIndexRef.current !== player.index) {
      sceneIndexRef.current = player.index;
      remainingRef.current = null;
    }

    if (!autoplayEnabled || isAnimating || !blockMotionDone || isPaused) {
      if (!autoplayEnabled || isAnimating || !blockMotionDone) {
        clearTimer();
      }
      return;
    }

    const delay = remainingRef.current ?? autoplayDelay;
    remainingRef.current = null;
    schedule(delay);

    return () => {
      if (deadlineRef.current !== null) {
        remainingRef.current = Math.max(0, deadlineRef.current - Date.now());
      }
      clearTimer();
    };
  }, [
    autoplayEnabled,
    autoplayDelay,
    isAnimating,
    blockMotionDone,
    isPaused,
    player.index,
    schedule,
    clearTimer,
  ]);
}
