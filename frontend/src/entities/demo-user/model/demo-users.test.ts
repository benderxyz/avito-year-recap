import { describe, expect, it } from 'vitest';
import { getDemoUserInitial } from '@/entities/demo-user';
import { mapUserDto } from '@/entities/demo-user/model/use-users-query';

describe('demo-user', () => {
  it('maps user dto to demo user', () => {
    expect(
      mapUserDto({
        user_id: 42,
        username: 'Alex',
        external_id: 'avito-42',
        timezone: 'UTC',
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }),
    ).toEqual({
      id: '42',
      name: 'Alex',
      externalId: 'avito-42',
    });
  });

  it('returns initial letter', () => {
    expect(getDemoUserInitial('Maria')).toBe('M');
  });
});
