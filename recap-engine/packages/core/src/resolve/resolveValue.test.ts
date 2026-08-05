import { describe, expect, it, jest } from '@jest/globals';

import { createFormatters } from '../format/createFormatters';
import { createTheme } from '../theme/createTheme';
import type { RecapContext } from '../types/context';
import { resolveValue } from './resolveValue';

const context: RecapContext<{ count: number }> = {
  data: { count: 4 },
  theme: createTheme(),
  index: 1,
  total: 3,
  format: createFormatters(),
};

describe('resolveValue', () => {
  it('returns undefined unchanged', () => {
    expect(resolveValue(undefined, context)).toBeUndefined();
  });

  it.each([
    ['text', 'text'],
    [0, 0],
    [false, false],
    [null, null],
  ])('returns literal value %p unchanged', (value, expected) => {
    expect(resolveValue(value, context)).toBe(expected);
  });

  it('calls a resolver exactly once with the full context', () => {
    const resolver = jest.fn((ctx: typeof context) => ctx.data.count + ctx.index);

    expect(resolveValue(resolver, context)).toBe(5);
    expect(resolver).toHaveBeenCalledTimes(1);
    expect(resolver).toHaveBeenCalledWith(context);
  });
});
