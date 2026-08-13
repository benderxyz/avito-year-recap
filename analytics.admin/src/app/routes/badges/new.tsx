import { createFileRoute } from '@tanstack/react-router';
import BadgePage from '@/pages/BadgePage';

export const Route = createFileRoute('/badges/new')({
  component: BadgePage,
});
