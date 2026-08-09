import { describe, expect, it } from 'vitest';
import { DEMO_USERS, getDemoUserById, getDemoUserInitial } from '@/entities/demo-user';

describe('demo-user', () => {
  it('contains 10 seed users', () => {
    expect(DEMO_USERS).toHaveLength(10);
  });

  it('finds user by id', () => {
    expect(getDemoUserById('42')?.name).toBe('Alex');
  });

  it('returns initial letter', () => {
    expect(getDemoUserInitial('Maria')).toBe('M');
  });
});
