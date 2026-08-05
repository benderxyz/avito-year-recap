import { describe, expect, it, jest } from '@jest/globals';

import type { Formatters, PluralForms } from '../types/context';
import { fillTemplate, resolveUnit } from './templates';

const forms: PluralForms = { one: 'товар', few: 'товара', many: 'товаров' };

describe('resolveUnit', () => {
  const plural = jest.fn<Formatters['plural']>(() => 'товара');
  const format = { plural } as unknown as Formatters;

  it('returns an empty string when unit is absent', () => {
    expect(resolveUnit(2, undefined, format)).toBe('');
    expect(plural).not.toHaveBeenCalled();
  });

  it('returns a string unit directly', () => {
    expect(resolveUnit(2, 'шт.', format)).toBe('шт.');
    expect(plural).not.toHaveBeenCalled();
  });

  it('delegates plural forms to the formatter', () => {
    expect(resolveUnit(2, forms, format)).toBe('товара');
    expect(plural).toHaveBeenCalledWith(2, forms);
  });
});

describe('fillTemplate', () => {
  it('replaces repeated word placeholders with string and number values', () => {
    expect(fillTemplate('{{name}}: {{count}} из {{count}}', { name: 'Продажи', count: 7 })).toBe(
      'Продажи: 7 из 7',
    );
  });

  it('replaces missing values with an empty string', () => {
    expect(fillTemplate('До {{missing}} после', {})).toBe('До  после');
  });

  it('leaves unsupported placeholder syntax untouched', () => {
    expect(fillTemplate('{{with-dash}} / {single}', { 'with-dash': 'x', single: 'y' })).toBe(
      '{{with-dash}} / {single}',
    );
  });
});
