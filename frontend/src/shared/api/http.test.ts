import { describe, expect, it } from 'vitest';
import { buildRecapUrl, buildShareRecapUrl } from '@/shared/api';

describe('shared api', () => {
  it('builds recap url without trailing slash', () => {
    expect(buildRecapUrl('http://localhost:8081/', 2026, '42')).toBe(
      'http://localhost:8081/api/recap/2026/42',
    );
    expect(buildRecapUrl('https://recaps.hakolr.dev', 2026, '42')).toBe(
      'https://recaps.hakolr.dev/api/recap/2026/42',
    );
  });

  it('builds share recap url without trailing slash', () => {
    expect(buildShareRecapUrl('http://localhost:8081/', 'abc.def')).toBe(
      'http://localhost:8081/api/share/abc.def',
    );
    expect(buildShareRecapUrl('https://recaps.hakolr.dev', 'NDk6MjAyNg.token')).toBe(
      'https://recaps.hakolr.dev/api/share/NDk6MjAyNg.token',
    );
  });
});
