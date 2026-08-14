import userEvent from '@testing-library/user-event';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useGetApiAdminRecommendations } from '@/entities/recommendation';
import type { RecommendationRule } from '@/shared/api/generated/model/recommendationRule';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import RecommendationsPage from './RecommendationsPage';

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

vi.mock('@/shared/api/generated/recommendations/recommendations', () => ({
  useGetApiAdminRecommendations: vi.fn(),
  usePutApiAdminRecommendationsId: () => ({ mutateAsync: vi.fn() }),
}));

const recommendation: RecommendationRule = {
  id: 'more_orders',
  title: 'More orders',
  text: 'Sell more',
  callout: 'Tip',
  linkLabel: 'Open',
  path: '/orders',
  enabled: true,
  priority: 10,
  when: {
    match: 'all',
    predicates: [{ metric: 'orders_count', op: 'gte', value: 1 }],
  },
  createdAt: '2026-01-01T00:00:00.000Z',
  updatedAt: '2026-01-01T00:00:00.000Z',
};

describe('RecommendationsPage', () => {
  beforeEach(() => {
    navigate.mockReset();
    vi.mocked(useGetApiAdminRecommendations).mockReturnValue({
      data: { items: [recommendation] },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useGetApiAdminRecommendations>);
  });

  it('renders recommendations and navigates on row click', async () => {
    const user = userEvent.setup();
    const { getByText, getByRole } = renderWithProviders(<RecommendationsPage />);

    expect(getByRole('link', { name: 'Create' })).toHaveAttribute('href', '/recommendations/new');
    expect(getByText('more_orders')).toBeInTheDocument();
    expect(getByText('More orders')).toBeInTheDocument();
    expect(getByText('/orders')).toBeInTheDocument();
    expect(getByText('orders_count')).toBeInTheDocument();

    await user.click(getByText('more_orders'));

    expect(navigate).toHaveBeenCalledWith({
      to: '/recommendations/$id',
      params: { id: 'more_orders' },
    });
  });

  it('reads filters from search params and writes them back', async () => {
    const user = userEvent.setup();
    const onUrlUpdate = vi.fn();
    const { getByLabelText } = renderWithProviders(<RecommendationsPage />, {
      searchParams: { search: 'more', enabled: 'true' },
      onUrlUpdate,
    });

    expect(getByLabelText('Search')).toHaveValue('more');

    await user.type(getByLabelText('Search'), 'x');

    const lastUpdate = onUrlUpdate.mock.calls.at(-1)?.[0] as
      | { searchParams: URLSearchParams }
      | undefined;
    expect(lastUpdate?.searchParams.get('search')).toBe('morex');
    expect(lastUpdate?.searchParams.get('enabled')).toBe('true');
  });
});
