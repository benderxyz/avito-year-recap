import { createFileRoute, redirect } from '@tanstack/react-router';
import { routes } from '@/shared/lib/routes';

export const Route = createFileRoute('/')({
  beforeLoad: () => {
    throw redirect({ to: routes.badges });
  },
});
