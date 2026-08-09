import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createMemoryHistory, createRouter, RouterProvider } from '@tanstack/react-router';
import { render, screen } from '@testing-library/react';
import { NuqsAdapter } from 'nuqs/adapters/tanstack-router';
import { describe, expect, it } from 'vitest';
import { routeTree } from '@/app/routeTree.gen';
import { mantineTheme } from '@/shared/ui';

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
  it('renders demo users list', async () => {
    const router = renderDemoAccountsRoute();
    await router.load();

    expect(screen.getByText('Demo аккаунты')).toBeInTheDocument();
    expect(screen.getByText('Alex')).toBeInTheDocument();
    expect(screen.getByText('Nina')).toBeInTheDocument();
    expect(screen.getAllByRole('link')).toHaveLength(10);
  });
});
