import { createRootRoute, Outlet } from '@tanstack/react-router';
import { NuqsAdapter } from 'nuqs/adapters/tanstack-router';
import { AppProviders } from '@/app/providers';

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <AppProviders>
      <NuqsAdapter>
        <Outlet />
      </NuqsAdapter>
    </AppProviders>
  );
}
