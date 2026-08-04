import { useEffect, useState } from 'react';

export function usePrefersReducedMotion(override?: boolean): boolean {
  const [prefers, setPrefers] = useState(false);

  useEffect(() => {
    if (override !== undefined) return;

    const mq = window.matchMedia('(prefers-reduced-motion: reduce)');

    const update = () => setPrefers(mq.matches);

    update();
    mq.addEventListener('change', update);

    return () => mq.removeEventListener('change', update);
  }, [override]);

  if (override !== undefined) return override;

  return prefers;
}
