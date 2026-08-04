import type { Formatters, PluralForms } from '../types/context';

export function resolveUnit(
  n: number,
  unit: string | PluralForms | undefined,
  format: Formatters,
): string {
  if (!unit) return '';

  if (typeof unit === 'string') return unit;

  return format.plural(n, unit);
}

export function fillTemplate(template: string, values: Record<string, string | number>): string {
  return template.replace(/\{\{(\w+)\}\}/g, (_, key: string) => String(values[key] ?? ''));
}
