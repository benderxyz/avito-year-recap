import type { Fn, RecapContext, ValueOrFn } from '../types/context';

export function resolveValue<TValue, TData>(
  value: ValueOrFn<TValue, TData> | undefined,
  ctx: RecapContext<TData>,
): TValue | undefined {
  if (value === undefined) return undefined;

  if (typeof value === 'function') {
    return (value as Fn<TValue, TData>)(ctx);
  }

  return value;
}
