import { createFormatters, EMotionPreset } from '@recap-engine/core';
import { useMemo, useState } from 'react';
import { RecapProvider } from '../context/RecapContext';
import { usePrefersReducedMotion } from '../hooks/usePrefersReducedMotion';
import { Background, Progress } from '../primitives';
import type { RecapProps } from './interface';
import { RecapShell } from './RecapShell';
import { useBlockMotion } from './useBlockMotion';
import { useRecapPlayer } from './useRecapPlayer';

export function Recap<TData>({
  theme,
  data,
  scenes,
  badges = [],
  motions = {
    default: EMotionPreset.SlideUp,
    back: EMotionPreset.SlideLeft,
    reducedMotion: EMotionPreset.Fade,
  },
  locale = 'ru-RU',
  initialSceneId,
  autoplay,
  loop,
  gestures = true,
  tapNav = false,
  holdToPause = false,
  reducedMotion: reducedMotionOverride,
  className,
  slots,
  onEvent,
}: RecapProps<TData>) {
  const reducedMotion = usePrefersReducedMotion(reducedMotionOverride);
  const format = useMemo(() => createFormatters(locale), [locale]);
  const [isAnimating, setAnimating] = useState(false);
  const [isAutoplayPaused, setAutoplayPaused] = useState(false);

  const { player, dispatch, next, prev, goTo, progress } = useRecapPlayer(
    scenes,
    initialSceneId,
    onEvent,
  );

  const { blockMotionDone, notifyBlockMotionComplete } = useBlockMotion(
    scenes[player.index],
    reducedMotion,
  );

  return (
    <RecapProvider
      value={{
        data,
        theme,
        scenes,
        player,
        dispatch,
        format,
        motions,
        reducedMotion,
        badges,
        onEvent,
        next,
        prev,
        goTo,
        progress,
        autoplay,
        isAutoplayPaused,
        isAnimating,
        setAnimating,
        notifyBlockMotionComplete,
        blockMotionDone,
      }}
    >
      <RecapShell
        gestures={gestures}
        tapNav={tapNav}
        holdToPause={holdToPause}
        autoplay={autoplay}
        loop={loop}
        isAutoplayPaused={isAutoplayPaused}
        onAutoplayPausedChange={setAutoplayPaused}
        slots={slots}
        className={className}
      />
    </RecapProvider>
  );
}

Recap.Progress = Progress;
Recap.Background = Background;
