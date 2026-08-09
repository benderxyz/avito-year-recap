import { useQueryStates } from 'nuqs';
import { productFilterParsers } from './parsers';

export function useProductFilters() {
  return useQueryStates(productFilterParsers, {
    history: 'replace',
    shallow: false,
  });
}
