import { describe, expect, it } from '@jest/globals';

import { EMetricType, type MetricValue } from '../types/payload';
import { metricDate, metricList, metricNumber, metricString } from './metrics';

const metrics: Record<string, MetricValue> = {
  number: { type: EMetricType.Number, value: 12 },
  money: { type: EMetricType.Money, value: 3400, currency: 'RUB' },
  percentile: { type: EMetricType.Percentile, value: 87 },
  ratio: { type: EMetricType.Ratio, value: 0.5 },
  string: { type: EMetricType.String, value: 'Москва' },
  date: { type: EMetricType.Date, value: '2024-03-14' },
  brokenDate: { type: EMetricType.Date, value: 'не дата' },
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

describe('metricDate', () => {
  it('parses an ISO date metric', () => {
    expect(metricDate(metrics, 'date')?.toISOString()).toBe('2024-03-14T00:00:00.000Z');
  });

  it('returns null for missing, non-date and unparsable metrics', () => {
    expect(metricDate(metrics, 'missing')).toBeNull();
    expect(metricDate(metrics, 'string')).toBeNull();
    expect(metricDate(metrics, 'brokenDate')).toBeNull();
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
