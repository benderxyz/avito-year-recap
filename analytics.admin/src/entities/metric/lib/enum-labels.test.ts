import { describe, expect, it } from 'vitest';
import { i18n } from '@/shared/i18n/config';
import { getMetricSourceFieldLabel, getMetricValueTypeLabel } from './enum-labels';

describe('metric enum labels', () => {
  it('localizes known valueType and sourceField values', () => {
    expect(getMetricValueTypeLabel(i18n.t.bind(i18n), 'number')).toBe('Number');
    expect(getMetricValueTypeLabel(i18n.t.bind(i18n), 'money')).toBe('Money');
    expect(getMetricSourceFieldLabel(i18n.t.bind(i18n), 'value')).toBe('Value');
    expect(getMetricSourceFieldLabel(i18n.t.bind(i18n), 'share')).toBe('Share');
  });
});
