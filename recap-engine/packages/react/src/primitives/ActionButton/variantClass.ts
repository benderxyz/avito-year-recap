import { EButtonVariant } from '@recap-engine/core';

export function variantClass(variant: EButtonVariant = EButtonVariant.Primary): string {
  return `recap-btn recap-btn--${variant}`;
}
