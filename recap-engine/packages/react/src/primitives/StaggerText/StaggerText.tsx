import { useRecap } from '../../context/RecapContext';
import { cn } from '../../utils/cn';
import type { StaggerTextProps } from './interface';

export function StaggerText({ text, className, animate }: StaggerTextProps) {
  const { reducedMotion, notifyBlockMotionComplete } = useRecap();
  const parts = text.split(/(?<=[.!?])\s+/).filter(Boolean);
  const shouldAnimate = animate && !reducedMotion;
  const items = parts.map((part, order) => ({ part, order }));

  return (
    <div
      className={cn('recap-stagger', className, shouldAnimate && 'recap-stagger--in')}
      onAnimationEnd={notifyBlockMotionComplete}
    >
      {items.map(({ part, order }) => (
        <p
          key={order}
          className="recap-stagger__item"
          style={{ animationDelay: shouldAnimate ? `${order * 100}ms` : undefined }}
        >
          {part}
        </p>
      ))}
    </div>
  );
}
