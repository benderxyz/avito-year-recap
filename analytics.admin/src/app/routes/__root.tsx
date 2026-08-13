import { createRootRoute } from '@tanstack/react-router';
import { NuqsAdapter } from 'nuqs/adapters/tanstack-router';
import AppLayout from '@/app/layout/AppLayout';
import AppProviders from '@/app/providers/AppProviders';

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <AppProviders>
      <NuqsAdapter>
        <AppLayout />
      </NuqsAdapter>
    </AppProviders>
  );
}
