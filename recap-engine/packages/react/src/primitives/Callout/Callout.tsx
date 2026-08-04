import { useEffect } from 'react';
import { useRecap } from '../../context/RecapContext';
import { cn } from '../../utils/cn';
import type { CalloutProps } from './interface';

export function Callout({ children, animate = true }: CalloutProps) {
  const { reducedMotion, notifyBlockMotionComplete } = useRecap();
  const shouldAnimate = animate && !reducedMotion;

  useEffect(() => {
    if (!children || shouldAnimate) return;
    notifyBlockMotionComplete();
  }, [children, shouldAnimate, notifyBlockMotionComplete]);

  if (!children) return null;

  return (
    <div
      className={cn('recap-callout', shouldAnimate && 'recap-callout--in')}
      onAnimationEnd={notifyBlockMotionComplete}
    >
      {children}
    </div>
  );
}
