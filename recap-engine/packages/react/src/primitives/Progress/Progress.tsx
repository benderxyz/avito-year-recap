import type { CSSProperties } from 'react';
import { useRecap } from '../../context/RecapContext';

function segmentFill(index: number, current: number, progress: number): number {
  if (index < current) return 1;
  if (index > current) return 0;
  return Math.max(progress, 0.08);
}

export function Progress() {
  const { progress, player } = useRecap();

  return (
    <div className="recap-progress" aria-hidden>
      {Array.from({ length: player.total }, (_, index) => (
        // Positional segments have no identity beyond index.
        // biome-ignore lint/suspicious/noArrayIndexKey: progress segments are positional
        <div key={index} className="recap-progress__segment">
          <div
            className="recap-progress__fill"
            style={
              {
                '--fill': segmentFill(index, player.index, progress),
              } as CSSProperties
            }
          />
        </div>
      ))}
    </div>
  );
}
