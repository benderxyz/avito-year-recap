import { parseAsBoolean, parseAsInteger, parseAsString } from 'nuqs';

export const recommendationFilterParsers = {
  search: parseAsString.withDefault(''),
  enabled: parseAsBoolean,
  metric: parseAsString.withDefault(''),
  minPriority: parseAsInteger,
};
