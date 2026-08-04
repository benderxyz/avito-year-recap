import { fillTemplate } from '@recap-engine/core';
import { cn } from '../../utils/cn';
import type { ComparisonProps } from './interface';

export function Comparison({ template, percentile, className }: ComparisonProps) {
  const text = fillTemplate(template, { percentile });
  return <p className={cn('recap-comparison', className)}>{text}</p>;
}
