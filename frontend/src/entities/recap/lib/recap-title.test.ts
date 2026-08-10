import { describe, expect, it } from 'vitest';
import { getPrivateRecapTitle, getSharedRecapTitle } from '@/entities/recap';

describe('recap title helpers', () => {
  it('builds private recap title', () => {
    expect(getPrivateRecapTitle(2026)).toBe('Итоги 2026');
    expect(getPrivateRecapTitle(undefined)).toBe('Итоги ');
  });

  it('builds shared recap title with display name', () => {
    expect(getSharedRecapTitle('Anna')).toBe('С вами поделилась Anna');
  });

  it('builds shared recap title fallback', () => {
    expect(getSharedRecapTitle(undefined)).toBe('С вами поделились итогами');
  });
});
