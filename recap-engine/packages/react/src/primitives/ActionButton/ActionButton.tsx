import { EButtonVariant } from '@recap-engine/core';
import { cn } from '../../utils/cn';
import type { ActionButtonProps } from './interface';
import { variantClass } from './variantClass';

export function ActionButton({
  variant = EButtonVariant.Primary,
  className,
  ...props
}: ActionButtonProps) {
  return <button type="button" className={cn(variantClass(variant), className)} {...props} />;
}
