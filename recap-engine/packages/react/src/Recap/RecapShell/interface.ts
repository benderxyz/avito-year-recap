import type { RecapProps } from '../interface';

export type RecapShellProps = {
  gestures?: boolean;
  tapNav?: RecapProps['tapNav'];
  holdToPause?: boolean;
  autoplay?: boolean | { delayMs: number };
  loop?: boolean;
  isAutoplayPaused?: boolean;
  onAutoplayPausedChange?: (paused: boolean) => void;
  slots?: RecapProps['slots'];
  className?: string;
};
