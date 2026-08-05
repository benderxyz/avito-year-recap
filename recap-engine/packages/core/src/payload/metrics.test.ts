import { describe, expect, it } from '@jest/globals';

import { EMetricType, type MetricValue } from '../types/payload';
import { metricList, metricNumber, metricString } from './metrics';

const metrics: Record<string, MetricValue> = {
  number: { type: EMetricType.Number, value: 12 },
  money: { type: EMetricType.Money, value: 3400, currency: 'RUB' },
  percentile: { type: EMetricType.Percentile, value: 87 },
  ratio: { type: EMetricType.Ratio, value: 0.5 },
  string: { type: EMetricType.String, value: 'Москва' },
  list: {
    type: EMetricType.List,
    value: [{ id: 'one', label: 'Первый', value: 1, imageUrl: '/one.png' }],
  },
};

describe('metricNumber', () => {
  it.each([
    ['number', 12],
    ['money', 3400],
    ['percentile', 87],
    ['ratio', 0.5],
  ])('reads numeric metric %s', (key, value) => {
    expect(metricNumber(metrics, key)).toBe(value);
  });

  it('returns the requested fallback for missing and non-numeric metrics', () => {
    expect(metricNumber(metrics, 'missing', 7)).toBe(7);
    expect(metricNumber(metrics, 'string', 8)).toBe(8);
    expect(metricNumber(metrics, 'list', 9)).toBe(9);
  });
});

describe('metricString', () => {
  it('reads string metrics', () => {
    expect(metricString(metrics, 'string')).toBe('Москва');
  });

  it('returns defaults for missing and non-string metrics', () => {
    expect(metricString(metrics, 'missing')).toBe('');
    expect(metricString(metrics, 'number', 'нет')).toBe('нет');
  });
});

describe('metricList', () => {
  it('returns the original list value', () => {
    expect(metricList(metrics, 'list')).toBe(metrics.list?.value);
  });

  it('returns an empty list for missing and non-list metrics', () => {
    expect(metricList(metrics, 'missing')).toEqual([]);
    expect(metricList(metrics, 'string')).toEqual([]);
  });
});
