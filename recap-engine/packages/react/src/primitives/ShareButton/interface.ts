import type { EButtonVariant, EShareKind } from '@recap-engine/core';

export type ShareButtonProps = {
  label: string;
  share: { kind: EShareKind; title?: string; text?: string; url?: string };
  variant?: EButtonVariant;
};
