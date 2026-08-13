import { describe, expect, it } from 'vitest';
import { metricFormDefaults, metricFormSchema } from './metric-form-schema';

describe('metricFormSchema', () => {
  it('rejects an empty key', () => {
    const result = metricFormSchema.safeParse(metricFormDefaults);

    expect(result.success).toBe(false);
  });

  it('accepts a valid metric', () => {
    const result = metricFormSchema.safeParse({
      ...metricFormDefaults,
      key: 'orders_count',
    });

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data.key).toBe('orders_count');
    }
  });
});
