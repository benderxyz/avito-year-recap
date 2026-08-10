import { describe, expect, it } from 'vitest';
import { resolveRecapColorScheme } from '@/entities/recap';

describe('resolveRecapColorScheme', () => {
  it('returns explicit light scheme', () => {
    expect(resolveRecapColorScheme('light', 'dark')).toBe('light');
  });

  it('returns explicit dark scheme', () => {
    expect(resolveRecapColorScheme('dark', 'light')).toBe('dark');
  });

  it('falls back to system scheme when auto', () => {
    expect(resolveRecapColorScheme('auto', 'dark')).toBe('dark');
    expect(resolveRecapColorScheme('auto', 'light')).toBe('light');
  });
});
