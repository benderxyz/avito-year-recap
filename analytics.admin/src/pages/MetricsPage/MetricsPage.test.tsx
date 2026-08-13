import userEvent from '@testing-library/user-event';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useGetApiAdminMetrics } from '@/shared/api/generated/metrics/metrics';
import type { MetricDefinition } from '@/shared/api/generated/model/metricDefinition';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import MetricsPage from './MetricsPage';

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

vi.mock('@/shared/api/generated/metrics/metrics', () => ({
  useGetApiAdminMetrics: vi.fn(),
}));

const metric: MetricDefinition = {
  key: 'orders_count',
  valueType: 'number',
  sourceField: 'value',
  sourceKey: 'orders',
  sortOrder: 1,
  enabled: true,
  isPublic: true,
  includeInLlm: false,
  createdAt: '2026-01-01T00:00:00.000Z',
  updatedAt: '2026-01-01T00:00:00.000Z',
};

describe('MetricsPage', () => {
  beforeEach(() => {
    navigate.mockReset();
    vi.mocked(useGetApiAdminMetrics).mockReturnValue({
      data: { items: [metric] },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useGetApiAdminMetrics>);
  });

  it('renders metrics and navigates on row click', async () => {
    const user = userEvent.setup();
    const { getByText, getByRole } = renderWithProviders(<MetricsPage />);

    expect(getByRole('link', { name: 'Create' })).toHaveAttribute('href', '/metrics/new');
    expect(getByText('orders_count')).toBeInTheDocument();
    expect(getByText('Number')).toBeInTheDocument();
    expect(getByText('Value')).toBeInTheDocument();

    await user.click(getByText('orders_count'));

    expect(navigate).toHaveBeenCalledWith({ to: '/metrics/$key', params: { key: 'orders_count' } });
  });

  it('reads filters from search params and writes them back', async () => {
    const user = userEvent.setup();
    const onUrlUpdate = vi.fn();
    const { getByLabelText } = renderWithProviders(<MetricsPage />, {
      searchParams: { search: 'orders', enabled: 'true' },
      onUrlUpdate,
    });

    expect(getByLabelText('Search')).toHaveValue('orders');

    await user.type(getByLabelText('Search'), 'x');

    const lastUpdate = onUrlUpdate.mock.calls.at(-1)?.[0] as
      | { searchParams: URLSearchParams }
      | undefined;
    expect(lastUpdate?.searchParams.get('search')).toBe('ordersx');
    expect(lastUpdate?.searchParams.get('enabled')).toBe('true');
  });
});
