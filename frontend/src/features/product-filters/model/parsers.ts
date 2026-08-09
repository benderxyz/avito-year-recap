import { parseAsInteger, parseAsString, parseAsStringEnum } from 'nuqs';

export const productFilterParsers = {
  q: parseAsString.withDefault(''),
  category: parseAsString.withDefault(''),
  sort: parseAsString.withDefault('title'),
  order: parseAsStringEnum(['asc', 'desc']).withDefault('asc'),
  priceMin: parseAsString.withDefault(''),
  priceMax: parseAsString.withDefault(''),
  page: parseAsInteger.withDefault(1),
};

export const PRODUCTS_PAGE_SIZE = 12;

export function getProductsSkip(page: number): number {
  return (page - 1) * PRODUCTS_PAGE_SIZE;
}

export function parsePriceFilter(value: string): number | null {
  if (!value.trim()) {
    return null;
  }
  const parsed = Number.parseFloat(value);
  return Number.isFinite(parsed) ? parsed : null;
}
