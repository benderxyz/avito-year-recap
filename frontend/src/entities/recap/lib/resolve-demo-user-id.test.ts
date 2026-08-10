import { describe, expect, it } from 'vitest';
import { resolveDemoUserId } from '@/entities/recap';

const users = [
  { id: '42', name: 'Anna', externalId: 'avito-42' },
  { id: '7', name: 'Bob', externalId: 'avito-7' },
];

describe('resolveDemoUserId', () => {
  it('returns numeric id as-is', () => {
    expect(resolveDemoUserId('42', users)).toBe('42');
  });

  it('maps external id to demo user id', () => {
    expect(resolveDemoUserId('avito-42', users)).toBe('42');
  });

  it('falls back to meta user id when users are unavailable', () => {
    expect(resolveDemoUserId('avito-42')).toBe('avito-42');
  });
});
