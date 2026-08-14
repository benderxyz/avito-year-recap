import { describe, expect, it } from 'vitest';
import {
  recommendationFormDefaults,
  recommendationFormSchema,
  toRecommendationCreate,
} from './form-schema';

const validRecommendation = {
  ...recommendationFormDefaults,
  id: 'more_orders',
  title: 'More orders',
  text: 'Sell more',
  linkLabel: 'Open',
  path: '/orders',
};

describe('recommendationFormSchema', () => {
  it('rejects an empty id', () => {
    const result = recommendationFormSchema.safeParse(recommendationFormDefaults);

    expect(result.success).toBe(false);
  });

  it('accepts empty predicates', () => {
    const result = recommendationFormSchema.safeParse(validRecommendation);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data.when.predicates).toEqual([]);
    }
  });

  it('writes GroupWhen on create', () => {
    const created = toRecommendationCreate({
      ...validRecommendation,
      when: {
        match: 'all',
        predicates: [{ metric: 'orders_count', op: 'gte', value: 10 }],
      },
    });

    expect(created.when).toEqual({
      match: 'all',
      predicates: [{ metric: 'orders_count', op: 'gte', value: 10 }],
    });
  });

  it('omits predicate value for exists', () => {
    const created = toRecommendationCreate({
      ...validRecommendation,
      when: {
        match: 'any',
        predicates: [{ metric: 'orders_count', op: 'exists', value: null }],
      },
    });

    expect(created.when).toEqual({
      match: 'any',
      predicates: [{ metric: 'orders_count', op: 'exists' }],
    });
  });
});
