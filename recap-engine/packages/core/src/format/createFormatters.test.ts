import { describe, expect, it } from '@jest/globals';

import type { PluralForms } from '../types/context';
import { createFormatters } from './createFormatters';

const forms: PluralForms = { one: 'товар', few: 'товара', many: 'товаров' };

describe('createFormatters', () => {
  it('formats numbers in the selected locale and forwards options', () => {
    const format = createFormatters('en-US');
    expect(format.number(12345.67)).toBe(new Intl.NumberFormat('en-US').format(12345.67));
    expect(format.number(0.256, { style: 'percent', maximumFractionDigits: 0 })).toBe('26%');
  });

  it('formats whole currency values with RUB by default', () => {
    const format = createFormatters('en-US');
    expect(format.currency(1234.6)).toBe(
      new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'RUB',
        maximumFractionDigits: 0,
      }).format(1234.6),
    );
  });

  it('supports a currency override', () => {
    const format = createFormatters('en-US');
    expect(format.currency(99, 'USD')).toBe('$99');
  });

  it.each([
    [1, 'товар'],
    [21, 'товар'],
    [-31, 'товар'],
    [2, 'товара'],
    [24, 'товара'],
    [-3, 'товара'],
    [0, 'товаров'],
    [5, 'товаров'],
    [11, 'товаров'],
    [14, 'товаров'],
    [111, 'товаров'],
  ])('selects the Russian plural form for %s', (count, expected) => {
    expect(createFormatters().plural(count as number, forms)).toBe(expected);
  });
});
