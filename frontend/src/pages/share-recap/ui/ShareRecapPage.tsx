import { useNavigate } from '@tanstack/react-router';
import { useUsersQuery } from '@/entities/demo-user';
import { resolveDemoUserId, useSharedRecapQuery } from '@/entities/recap';
import { RecapModal } from '@/widgets/recap-modal';

type ShareRecapPageProps = {
  token: string;
};

export function ShareRecapPage({ token }: ShareRecapPageProps) {
  const navigate = useNavigate();
  const { data } = useSharedRecapQuery(token, true);
  const { data: users } = useUsersQuery();

  const handleClose = () => {
    if (data) {
      const authorId = resolveDemoUserId(data.meta.user.id, users);
      navigate({ to: '/demo/$id', params: { id: authorId } });
      return;
    }

    navigate({ to: '/demo' });
  };

  return (
    <div className="page-shell">
      <RecapModal mode="shared" shareToken={token} opened onClose={handleClose} />
    </div>
  );
}
