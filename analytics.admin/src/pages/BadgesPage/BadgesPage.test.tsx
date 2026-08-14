import userEvent from '@testing-library/user-event';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useGetApiAdminBadges } from '@/entities/badge';
import type { BadgeRule } from '@/shared/api/generated/model/badgeRule';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import BadgesPage from './BadgesPage';

const navigate = vi.fn();

vi.mock('@tanstack/react-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@tanstack/react-router')>();

  return {
    ...actual,
    useNavigate: () => navigate,
    Link: ({ to, children, ...props }: { to: string; children?: ReactNode }) => (
      <a href={typeof to === 'string' ? to : ''} {...props}>
        {children}
      </a>
    ),
  };
});

vi.mock('@/shared/api/generated/badges/badges', () => ({
  useGetApiAdminBadges: vi.fn(),
  usePutApiAdminBadgesId: () => ({ mutateAsync: vi.fn() }),
}));

const badge: BadgeRule = {
  id: 'top_seller',
  title: 'Top seller',
  description: 'Sold a lot',
  enabled: true,
  sortOrder: 1,
  visibility: 'public',
  when: { metric: 'orders_count', op: 'gte', value: 10 },
  createdAt: '2026-01-01T00:00:00.000Z',
  updatedAt: '2026-01-01T00:00:00.000Z',
};

describe('BadgesPage', () => {
  beforeEach(() => {
    navigate.mockReset();
    vi.mocked(useGetApiAdminBadges).mockReturnValue({
      data: { items: [badge] },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useGetApiAdminBadges>);
  });

  it('renders badges and navigates on row click', async () => {
    const user = userEvent.setup();
    const { getByText, getByRole } = renderWithProviders(<BadgesPage />);

    expect(getByRole('link', { name: 'Create' })).toHaveAttribute('href', '/badges/new');
    expect(getByText('top_seller')).toBeInTheDocument();
    expect(getByText('Top seller')).toBeInTheDocument();
    expect(getByRole('cell', { name: 'Public' })).toBeInTheDocument();
    expect(getByText('orders_count')).toBeInTheDocument();

    await user.click(getByText('top_seller'));

    expect(navigate).toHaveBeenCalledWith({ to: '/badges/$id', params: { id: 'top_seller' } });
  });

  it('reads filters from search params and writes them back', async () => {
    const user = userEvent.setup();
    const onUrlUpdate = vi.fn();
    const { getByLabelText } = renderWithProviders(<BadgesPage />, {
      searchParams: { search: 'top', enabled: 'true' },
      onUrlUpdate,
    });

    expect(getByLabelText('Search')).toHaveValue('top');

    await user.type(getByLabelText('Search'), 'x');

    const lastUpdate = onUrlUpdate.mock.calls.at(-1)?.[0] as
      | { searchParams: URLSearchParams }
      | undefined;
    expect(lastUpdate?.searchParams.get('search')).toBe('topx');
    expect(lastUpdate?.searchParams.get('enabled')).toBe('true');
  });
});
