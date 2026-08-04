import type { EButtonVariant } from '@recap-engine/core';
import type { ButtonHTMLAttributes } from 'react';

export type ActionButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: EButtonVariant;
};
