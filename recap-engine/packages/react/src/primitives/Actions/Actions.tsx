import type { ActionsProps } from './interface';

export function Actions({ children }: ActionsProps) {
  return <div className="recap-actions">{children}</div>;
}
