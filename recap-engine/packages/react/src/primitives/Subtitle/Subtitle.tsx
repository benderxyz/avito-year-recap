import { cn } from '../../utils/cn';
import type { SubtitleProps } from './interface';

export function Subtitle({ children, className }: SubtitleProps) {
  if (!children) return null;
  return <p className={cn('recap-subtitle', className)}>{children}</p>;
}
