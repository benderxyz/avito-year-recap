import { describe, expect, it } from 'vitest';
import { i18n } from '@/shared/i18n/config';
import { getMatchModeLabel, getPredicateOpLabel } from './enum-labels';

describe('recommendation enum labels', () => {
  it('localizes known match and predicate op values', () => {
    expect(getMatchModeLabel(i18n.t.bind(i18n), 'all')).toBe('All');
    expect(getMatchModeLabel(i18n.t.bind(i18n), 'any')).toBe('Any');
    expect(getPredicateOpLabel(i18n.t.bind(i18n), 'gt')).toBe('Greater than');
    expect(getPredicateOpLabel(i18n.t.bind(i18n), 'gte')).toBe('Greater or equal');
    expect(getPredicateOpLabel(i18n.t.bind(i18n), 'eq')).toBe('Equal');
    expect(getPredicateOpLabel(i18n.t.bind(i18n), 'exists')).toBe('Exists');
  });
});
