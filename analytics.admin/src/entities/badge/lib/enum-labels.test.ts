import { describe, expect, it } from 'vitest';
import { i18n } from '@/shared/i18n/config';
import { getBadgeVisibilityLabel, getPredicateOpLabel } from './enum-labels';

describe('badge enum labels', () => {
  it('localizes known visibility and predicateOp values', () => {
    expect(getBadgeVisibilityLabel(i18n.t.bind(i18n), 'private')).toBe('Private');
    expect(getBadgeVisibilityLabel(i18n.t.bind(i18n), 'public')).toBe('Public');
    expect(getPredicateOpLabel(i18n.t.bind(i18n), 'gt')).toBe('Greater than');
    expect(getPredicateOpLabel(i18n.t.bind(i18n), 'exists')).toBe('Exists');
  });
});
