import type { BackgroundProps } from './interface';

export function Background({ children }: BackgroundProps) {
  return <div className="recap-background">{children}</div>;
}
