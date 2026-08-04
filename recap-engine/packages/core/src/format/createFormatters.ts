import type { Formatters, PluralForms } from '../types/context';

function pluralRu(n: number, forms: PluralForms): string {
  const abs = Math.abs(n) % 100;
  const last = abs % 10;

  if (abs > 10 && abs < 20) return forms.many;
  if (last === 1) return forms.one;
  if (last >= 2 && last <= 4) return forms.few;

  return forms.many;
}

export function createFormatters(locale = 'ru-RU'): Formatters {
  return {
    number: (count, opts) => new Intl.NumberFormat(locale, opts).format(count),
    currency: (count, currency = 'RUB') =>
      new Intl.NumberFormat(locale, {
        style: 'currency',
        currency,
        maximumFractionDigits: 0,
      }).format(count),
    plural: (count, forms) => pluralRu(count, forms),
  };
}
