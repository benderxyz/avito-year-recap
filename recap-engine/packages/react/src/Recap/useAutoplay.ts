import { ERecapEventType } from '@recap-engine/core';
import { useEffect, useRef } from 'react';
import { useRecap } from '../context/RecapContext';

export function useAutoplay(options: {
  autoplay?: boolean | { delayMs: number };
  loop?: boolean;
  durationMs?: number;
}) {
  const { autoplay, loop = false, durationMs } = options;
  const { player, next, goTo, onEvent, isAnimating, blockMotionDone } = useRecap();
  const autoplayTimer = useRef<number | null>(null);

  const autoplayEnabled = Boolean(autoplay);
  const autoplayDelay = typeof autoplay === 'object' ? autoplay.delayMs : (durationMs ?? 3200);

  useEffect(() => {
    if (!autoplayEnabled || isAnimating || !blockMotionDone) return;
    if (autoplayTimer.current) window.clearTimeout(autoplayTimer.current);

    autoplayTimer.current = window.setTimeout(() => {
      if (player.index >= player.total - 1) {
        onEvent?.({ type: ERecapEventType.Complete });
        if (loop) goTo(0);
        return;
      }
      next();
    }, autoplayDelay);

    return () => {
      if (autoplayTimer.current) window.clearTimeout(autoplayTimer.current);
    };
  }, [
    autoplayEnabled,
    autoplayDelay,
    isAnimating,
    blockMotionDone,
    player.index,
    player.total,
    next,
    goTo,
    loop,
    onEvent,
  ]);
}
