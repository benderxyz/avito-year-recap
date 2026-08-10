import { createFileRoute } from '@tanstack/react-router';
import { DemoAccountsPage } from '@/pages/demo-accounts';

export const Route = createFileRoute('/demo/')({
  component: DemoAccountsPage,
});
