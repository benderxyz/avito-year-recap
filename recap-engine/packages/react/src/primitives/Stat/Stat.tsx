import { resolveUnit } from '@recap-engine/core';
import { useCallback } from 'react';
import { useRecap } from '../../context/RecapContext';
import { useCountUp } from '../../hooks/useCountUp';
import { cn } from '../../utils/cn';
import type { StatProps } from './interface';

export function Stat({
  value,
  unit,
  valueFormat,
  animate = true,
  classNames,
  onMotionComplete,
}: StatProps) {
  const { format, reducedMotion, notifyBlockMotionComplete } = useRecap();
  const handleComplete = useCallback(() => {
    onMotionComplete?.();
    notifyBlockMotionComplete();
  }, [onMotionComplete, notifyBlockMotionComplete]);

  const display = useCountUp(value, animate && !reducedMotion, 1000, handleComplete);
  const unitText = resolveUnit(value, unit, format);

  return (
    <div className={cn('recap-stat', classNames?.root)}>
      <div className={cn('recap-stat__value', classNames?.value)}>
        {format.number(display, valueFormat)}
      </div>
      {unitText ? <div className={cn('recap-stat__unit', classNames?.unit)}>{unitText}</div> : null}
    </div>
  );
}
