import { createFileRoute } from '@tanstack/react-router';
import { DemoCatalogPage } from '@/pages/demo-catalog';

export const Route = createFileRoute('/demo/$id')({
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = Route.useParams();
  return <DemoCatalogPage userId={id} />;
}
