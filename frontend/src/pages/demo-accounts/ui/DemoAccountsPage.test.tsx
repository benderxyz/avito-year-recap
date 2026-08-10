import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createMemoryHistory, createRouter, RouterProvider } from '@tanstack/react-router';
import { render, screen, waitFor } from '@testing-library/react';
import { NuqsAdapter } from 'nuqs/adapters/tanstack-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { routeTree } from '@/app/routeTree.gen';
import { fetchUsers } from '@/shared/api/users-api';
import { mantineTheme } from '@/shared/ui';

vi.mock('@/shared/api/users-api', () => ({
  fetchUsers: vi.fn(),
}));

const mockUsers = [
  {
    user_id: 42,
    external_id: 'avito-42',
    username: 'Alex',
    timezone: 'UTC',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
  },
  {
    user_id: 51,
    external_id: 'avito-51',
    username: 'Nina',
    timezone: 'UTC',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
  },
];

function renderDemoAccountsRoute() {
  const history = createMemoryHistory({ initialEntries: ['/demo'] });
  const router = createRouter({ routeTree, history });
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  render(
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={mantineTheme}>
        <NuqsAdapter>
          <RouterProvider router={router} />
        </NuqsAdapter>
      </MantineProvider>
    </QueryClientProvider>,
  );

  return router;
}

describe('DemoAccountsPage', () => {
  beforeEach(() => {
    vi.mocked(fetchUsers).mockResolvedValue(mockUsers);
  });

  it('renders demo users list', async () => {
    const router = renderDemoAccountsRoute();
    await router.load();

    expect(screen.getByText('Demo аккаунты')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText('Alex')).toBeInTheDocument();
      expect(screen.getByText('Nina')).toBeInTheDocument();
      expect(screen.getAllByRole('link')).toHaveLength(2);
    });
  });
});
