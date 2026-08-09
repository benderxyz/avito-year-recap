import { describe, expect, it } from 'vitest';
import { buildRecapUrl, buildUsersListUrl } from '@/shared/api';

describe('shared api', () => {
  it('builds recap url without trailing slash', () => {
    expect(buildRecapUrl('http://localhost:8081/', 2026, '42')).toBe(
      'http://localhost:8081/api/recap/2026/42',
    );
  });

  it('builds users list url without trailing slash', () => {
    expect(buildUsersListUrl('http://localhost:8082/')).toBe('http://localhost:8082/users');
    expect(buildUsersListUrl('https://recap.hakolr.dev/api')).toBe(
      'https://recap.hakolr.dev/api/users',
    );
  });
});
