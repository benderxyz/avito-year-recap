import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createMemoryHistory, createRouter, RouterProvider } from '@tanstack/react-router';
import { type RenderOptions, render } from '@testing-library/react';
import { NuqsAdapter } from 'nuqs/adapters/tanstack-router';
import type { ReactElement, ReactNode } from 'react';
import { routeTree } from '@/app/routeTree.gen';
import { mantineTheme } from '@/shared/ui';

export function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });
}

export function renderWithProviders(ui: ReactElement, options?: RenderOptions) {
  const queryClient = createTestQueryClient();

  function Wrapper({ children }: { children: ReactNode }) {
    const history = createMemoryHistory({ initialEntries: ['/demo'] });
    const router = createRouter({
      routeTree,
      history,
      context: { queryClient },
    });

    return (
      <QueryClientProvider client={queryClient}>
        <MantineProvider theme={mantineTheme}>
          <RouterProvider router={router} />
          <NuqsAdapter>{children}</NuqsAdapter>
        </MantineProvider>
      </QueryClientProvider>
    );
  }

  return render(ui, { wrapper: Wrapper, ...options });
}
