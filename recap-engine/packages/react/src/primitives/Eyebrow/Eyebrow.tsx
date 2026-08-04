import { cn } from '../../utils/cn';
import type { EyebrowProps } from './interface';

export function Eyebrow({ children, className }: EyebrowProps) {
  if (!children) return null;
  return <p className={cn('recap-eyebrow', className)}>{children}</p>;
}
