import { createFileRoute } from '@tanstack/react-router';
import { ShareRecapPage } from '@/pages/share-recap';

export const Route = createFileRoute('/share/$token')({
  component: RouteComponent,
});

function RouteComponent() {
  const { token } = Route.useParams();
  return <ShareRecapPage token={token} />;
}
