import { EMetricType, type MetricListItem, type MetricValue } from '../types/payload';

export function metricNumber(
  metrics: Record<string, MetricValue>,
  key: string,
  fallback = 0,
): number {
  const m = metrics[key];

  if (!m) return fallback;

  if (
    m.type === EMetricType.Number ||
    m.type === EMetricType.Money ||
    m.type === EMetricType.Percentile ||
    m.type === EMetricType.Ratio
  ) {
    return m.value;
  }

  return fallback;
}

export function metricString(
  metrics: Record<string, MetricValue>,
  key: string,
  fallback = '',
): string {
  const m = metrics[key];

  if (!m) return fallback;

  if (m.type === EMetricType.String) return m.value;

  return fallback;
}

export function metricDate(metrics: Record<string, MetricValue>, key: string): Date | null {
  const m = metrics[key];

  if (m?.type !== EMetricType.Date) return null;

  const parsed = new Date(m.value);

  if (Number.isNaN(parsed.getTime())) return null;

  return parsed;
}

export function metricList(metrics: Record<string, MetricValue>, key: string): MetricListItem[] {
  const m = metrics[key];

  if (m?.type !== EMetricType.List) return [];

  return m.value;
}
