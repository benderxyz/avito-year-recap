import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiError, customFetch } from '@/shared/api/http';

afterEach(() => {
  vi.unstubAllEnvs();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe('customFetch', () => {
  it('prefixes base url and sends bearer token', async () => {
    vi.stubEnv('VITE_API_BASE_URL', 'https://api.example.com');
    vi.stubEnv('VITE_ADMIN_API_TOKEN', 'secret');

    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ items: [] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await expect(customFetch('/api/admin/badges')).resolves.toEqual({ items: [] });

    expect(fetchMock).toHaveBeenCalledOnce();
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toBe('https://api.example.com/api/admin/badges');
    expect((init.headers as Headers).get('Authorization')).toBe('Bearer secret');
  });

  it('returns undefined for 204', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(null, { status: 204 })));

    await expect(customFetch('/api/admin/badges/1', { method: 'DELETE' })).resolves.toBeUndefined();
  });

  it('throws ApiError with backend message', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ error: 'Нет токена или токен неверный' }), {
          status: 401,
          statusText: 'Unauthorized',
          headers: { 'Content-Type': 'application/json' },
        }),
      ),
    );

    const error = await customFetch('/api/admin/badges').catch((reason: unknown) => reason);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      message: 'Нет токена или токен неверный',
      status: 401,
    });
  });
});
