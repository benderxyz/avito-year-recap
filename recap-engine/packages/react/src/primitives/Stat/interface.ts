import type { PluralForms } from '@recap-engine/core';

export type StatProps = {
  value: number;
  unit?: string | PluralForms;
  valueFormat?: Intl.NumberFormatOptions;
  animate?: boolean;
  classNames?: { root?: string; value?: string; unit?: string };
  onMotionComplete?: () => void;
};
