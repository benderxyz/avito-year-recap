import { cn } from '../../utils/cn';
import type { TitleProps } from './interface';

export function Title({ children, className }: TitleProps) {
  if (!children) return null;
  return <h2 className={cn('recap-title', className)}>{children}</h2>;
}
