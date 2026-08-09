import { describe, expect, it } from 'vitest';
import { buildRecapUrl } from '@/shared/api';

describe('shared api', () => {
  it('builds recap url without trailing slash', () => {
    expect(buildRecapUrl('http://localhost:8081/', 2026, '42')).toBe(
      'http://localhost:8081/api/recap/2026/42',
    );
  });
});
