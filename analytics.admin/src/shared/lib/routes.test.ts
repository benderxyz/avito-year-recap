import { describe, expect, it } from 'vitest';
import { routes } from '@/shared/lib/routes';

describe('routes', () => {
  it('exposes absolute admin paths', () => {
    expect(routes.badges).toBe('/badges');
    expect(routes.badgeById).toBe('/badges/$id');
    expect(routes.metrics).toBe('/metrics');
    expect(routes.metricByKey).toBe('/metrics/$key');
    expect(routes.metricCreate).toBe('/metrics/new');
    expect(routes.preview).toBe('/preview');
  });
});
