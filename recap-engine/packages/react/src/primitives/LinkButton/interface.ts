import type { EButtonVariant, ELinkTarget } from '@recap-engine/core';

export type LinkButtonProps = {
  label: string;
  href: string;
  target?: ELinkTarget;
  variant?: EButtonVariant;
};
