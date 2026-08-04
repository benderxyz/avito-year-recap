import { ERecapEventType } from '@recap-engine/core';
import { useCallback, useEffect } from 'react';
import { useRecap } from '../context/RecapContext';

export function useKeyboardNav() {
  const { next, prev, player, onEvent } = useRecap();

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'ArrowRight' || e.key === ' ') {
        e.preventDefault();
        if (player.index >= player.total - 1) {
          onEvent?.({ type: ERecapEventType.Complete });
        }
        next();
      }
      if (e.key === 'ArrowLeft') {
        e.preventDefault();
        prev();
      }
    },
    [next, onEvent, player.index, player.total, prev],
  );

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleKeyDown]);
}
