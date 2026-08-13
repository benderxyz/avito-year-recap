import { describe, expect, it } from 'vitest';
import { badgeFormDefaults, badgeFormSchema, toBadgeCreate } from './form-schema';

const validBadge = {
  ...badgeFormDefaults,
  id: 'top_seller',
  title: 'Top seller',
  description: 'Sold a lot',
  when: {
    metric: 'orders_count',
    op: 'gte' as const,
    value: 10,
  },
};

describe('badgeFormSchema', () => {
  it('rejects an empty id', () => {
    const result = badgeFormSchema.safeParse(badgeFormDefaults);

    expect(result.success).toBe(false);
  });

  it('accepts a valid badge', () => {
    const result = badgeFormSchema.safeParse(validBadge);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data.id).toBe('top_seller');
    }
  });

  it('omits when.value for exists', () => {
    const created = toBadgeCreate({
      ...validBadge,
      when: { metric: 'orders_count', op: 'exists', value: null },
    });

    expect(created.when).toEqual({ metric: 'orders_count', op: 'exists' });
  });
});
