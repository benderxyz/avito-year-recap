import type { CSSProperties } from 'react';
import { useRecap } from '../../context/RecapContext';
import { cn } from '../../utils/cn';
import type { AchievementBadgeProps } from './interface';

export function AchievementBadge({
  title,
  description,
  icon,
  animate = true,
}: AchievementBadgeProps) {
  const { reducedMotion, notifyBlockMotionComplete } = useRecap();

  const glyphStyle = icon
    ? ({
        WebkitMaskImage: `url(${icon})`,
        maskImage: `url(${icon})`,
      } satisfies CSSProperties)
    : undefined;

  return (
    <div
      className={cn('recap-achievement', animate && !reducedMotion && 'recap-achievement--pop')}
      onAnimationEnd={notifyBlockMotionComplete}
    >
      {icon ? (
        <div className="recap-achievement__icon" aria-hidden>
          <span className="recap-achievement__icon-glyph" style={glyphStyle} />
        </div>
      ) : null}
      <h2 className="recap-achievement__title">{title}</h2>
      {description ? <p className="recap-achievement__desc">{description}</p> : null}
    </div>
  );
}
