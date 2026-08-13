import { describe, expect, it } from 'vitest';
import { createQueryClient } from '@/shared/api/query-client';

describe('createQueryClient', () => {
  it('disables refetch on window focus and retries once', () => {
    const client = createQueryClient();

    expect(client.getDefaultOptions().queries?.refetchOnWindowFocus).toBe(false);
    expect(client.getDefaultOptions().queries?.retry).toBe(1);
  });
});
