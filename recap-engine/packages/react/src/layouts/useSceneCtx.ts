import type { RecapContext } from '@recap-engine/core';
import { useRecap } from '../context/RecapContext';

export function useSceneCtx<TData>(): RecapContext<TData> & {
  goNext: () => void;
  goPrev: () => void;
  reducedMotion: boolean;
  badges: { id: string; title: string; description: string; icon?: string }[];
  notifyBlockMotionComplete: () => void;
} {
  const ctx = useRecap<TData>();

  return {
    data: ctx.data,
    theme: ctx.theme,
    index: ctx.player.index,
    total: ctx.player.total,
    format: ctx.format,
    goNext: ctx.next,
    goPrev: ctx.prev,
    reducedMotion: ctx.reducedMotion,
    badges: ctx.badges,
    notifyBlockMotionComplete: ctx.notifyBlockMotionComplete,
  };
}
