import { createFileRoute } from '@tanstack/react-router';
import MetricPage from '@/pages/MetricPage';

export const Route = createFileRoute('/metrics/new')({
  component: MetricPage,
});
