import { describe, expect, it } from 'vitest';
import {
  getProductsSkip,
  parsePriceFilter,
  productFilterParsers,
} from '@/features/product-filters';

describe('product filter parsers', () => {
  it('has default values', () => {
    expect(productFilterParsers.q.defaultValue).toBe('');
    expect(productFilterParsers.sort.defaultValue).toBe('title');
    expect(productFilterParsers.page.defaultValue).toBe(1);
  });

  it('calculates skip from page', () => {
    expect(getProductsSkip(1)).toBe(0);
    expect(getProductsSkip(3)).toBe(24);
  });

  it('parses price filter values', () => {
    expect(parsePriceFilter('')).toBeNull();
    expect(parsePriceFilter('100')).toBe(100);
    expect(parsePriceFilter('abc')).toBeNull();
  });
});
