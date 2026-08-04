import { useEffect, useState } from 'react';

export function useCountUp(
  target: number,
  enabled: boolean,
  durationMs = 1000,
  onComplete?: () => void,
): number {
  const [value, setValue] = useState(enabled ? 0 : target);

  useEffect(() => {
    if (!enabled) {
      setValue(target);
      onComplete?.();
      return;
    }

    setValue(0);
    let frame = 0;
    const start = performance.now();

    const tick = (now: number) => {
      const t = Math.min(1, (now - start) / durationMs);
      const eased = 1 - (1 - t) ** 3;

      setValue(Math.round(target * eased));

      if (t < 1) {
        frame = requestAnimationFrame(tick);
      } else {
        setValue(target);
        onComplete?.();
      }
    };

    frame = requestAnimationFrame(tick);

    return () => cancelAnimationFrame(frame);
  }, [target, enabled, durationMs, onComplete]);

  return value;
}
