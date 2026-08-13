import type { CSSProperties } from 'react';
import { useRecap } from '../../context/RecapContext';

const DEFAULT_AUTOPLAY_DELAY_MS = 3200;

export function Progress() {
  const {
    autoplay,
    blockMotionDone,
    goTo,
    isAnimating,
    isAutoplayPaused,
    player,
    reducedMotion,
    scenes,
  } = useRecap();

  if (player.total <= 1) return null;

  const currentScene = scenes[player.index];
  const durationMs =
    typeof autoplay === 'object'
      ? autoplay.delayMs
      : (currentScene?.durationMs ?? DEFAULT_AUTOPLAY_DELAY_MS);
  const animated = Boolean(autoplay) && !reducedMotion;
  const playing = animated && !isAnimating && blockMotionDone && !isAutoplayPaused;

  return (
    <nav className="recap-progress" aria-label="Навигация по сценам">
      {Array.from({ length: player.total }, (_, index) => (
        <button
          key={scenes[index]?.id ?? index}
          type="button"
          className="recap-progress__segment"
          aria-label={`Перейти к сцене ${index + 1} из ${player.total}`}
          aria-current={index === player.index ? 'step' : undefined}
          onClick={() => goTo(index)}
        >
          <span className="recap-progress__track">
            <span
              className="recap-progress__fill"
              data-state={
                index < player.index ? 'complete' : index === player.index ? 'current' : 'pending'
              }
              data-animated={index === player.index && animated}
              data-playing={index === player.index && playing}
              style={
                {
                  '--recap-progress-duration': `${durationMs}ms`,
                } as CSSProperties
              }
            />
          </span>
        </button>
      ))}
    </nav>
  );
}
