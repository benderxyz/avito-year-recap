import type { RecapProps } from '../interface';

export type RecapShellProps = {
  gestures?: boolean;
  autoplay?: boolean | { delayMs: number };
  loop?: boolean;
  slots?: RecapProps['slots'];
  className?: string;
};
