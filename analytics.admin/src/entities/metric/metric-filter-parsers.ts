import { parseAsBoolean, parseAsString } from 'nuqs';

export const metricFilterParsers = {
  search: parseAsString.withDefault(''),
  enabled: parseAsBoolean,
  isPublic: parseAsBoolean,
  includeInLlm: parseAsBoolean,
};
