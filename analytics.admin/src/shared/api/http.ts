export class ApiError extends Error {
  readonly status: number;
  readonly body?: unknown;

  constructor(message: string, status: number, body?: unknown) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.body = body;
  }
}

export async function customFetch<T>(url: string, options?: RequestInit): Promise<T> {
  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';
  const token = import.meta.env.VITE_ADMIN_API_TOKEN;
  const headers = new Headers(options?.headers);

  headers.set('Accept', 'application/json');

  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  if (options?.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(`${baseUrl}${url}`, {
    ...options,
    headers,
  });

  if (response.status === 204 || response.status === 205) {
    return undefined as T;
  }

  const isJson = response.headers.get('content-type')?.includes('application/json');
  const body: unknown = isJson ? await response.json() : await response.text();

  if (!response.ok) {
    const message =
      isJson && typeof body === 'object' && body !== null && 'error' in body
        ? String((body as { error: string }).error)
        : response.statusText;
    throw new ApiError(message, response.status, body);
  }

  return body as T;
}
