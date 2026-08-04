import type { CSSProperties } from 'react';
import { useRecap } from '../../context/RecapContext';
import { useSwipe } from '../../hooks/useSwipe';
import { SceneRenderer } from '../../layouts';
import { Background, Progress } from '../../primitives';
import { cn } from '../../utils/cn';
import { resolveSceneBackground } from '../resolveSceneBackground';
import { useAutoplay } from '../useAutoplay';
import { useKeyboardNav } from '../useKeyboardNav';
import { useSceneLifecycle } from '../useSceneLifecycle';
import { useScenePresentation } from '../useScenePresentation';
import type { RecapShellProps } from './interface';

export function RecapShell<TData>({
  gestures = true,
  autoplay,
  loop = false,
  slots,
  className,
}: RecapShellProps) {
  const {
    theme,
    scenes,
    player,
    next,
    prev,
    onEvent,
    reducedMotion,
    motions,
    setAnimating,
    dispatch,
  } = useRecap<TData>();

  const scene = scenes[player.index];
  const { motion, motionClass } = useScenePresentation(
    scene,
    player.direction,
    motions,
    reducedMotion,
  );

  useSceneLifecycle({
    scene,
    scenes,
    index: player.index,
    durationMs: motion.durationMs,
    reducedMotion,
    onEvent,
    setAnimating,
    dispatch,
  });

  useAutoplay({ autoplay, loop, durationMs: scene?.durationMs });
  useKeyboardNav();

  const swipe = useSwipe(gestures, next, prev);

  if (!scene) return null;

  const background = resolveSceneBackground(scene, theme.assets.background);

  return (
    <div className={cn('recap-root', className)} style={theme.cssVars as CSSProperties} {...swipe}>
      <Background>
        {background.useSlotBg ? slots?.background : null}
        <div className="recap-background__base" style={background.style} />
        <div className="recap-background__pattern" />
      </Background>

      {slots?.header ?? <Progress />}

      <section
        key={scene.id}
        className={cn('recap-scene', motionClass)}
        style={
          {
            '--recap-motion-duration': `${motion.durationMs}ms`,
          } as CSSProperties
        }
        aria-roledescription="slide"
        aria-label={`Сцена ${player.index + 1} из ${player.total}`}
      >
        <SceneRenderer scene={scene} />
      </section>

      {slots?.footer}
    </div>
  );
}
