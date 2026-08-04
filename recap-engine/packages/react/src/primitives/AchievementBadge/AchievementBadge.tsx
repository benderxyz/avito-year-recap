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

  return (
    <div
      className={cn('recap-achievement', animate && !reducedMotion && 'recap-achievement--pop')}
      onAnimationEnd={notifyBlockMotionComplete}
    >
      <div className="recap-achievement__icon" aria-hidden>
        {icon ?? '★'}
      </div>
      <h2 className="recap-achievement__title">{title}</h2>
      {description ? <p className="recap-achievement__desc">{description}</p> : null}
    </div>
  );
}
