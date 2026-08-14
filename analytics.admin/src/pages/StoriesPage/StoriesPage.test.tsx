import userEvent from '@testing-library/user-event';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useGetApiAdminStories } from '@/entities/story';
import type { StoryRule } from '@/shared/api/generated/model/storyRule';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import StoriesPage from './StoriesPage';

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

vi.mock('@/shared/api/generated/stories/stories', () => ({
  useGetApiAdminStories: vi.fn(),
  usePutApiAdminStoriesId: () => ({ mutateAsync: vi.fn() }),
}));

const story: StoryRule = {
  id: 'year_intro',
  sceneType: 'intro',
  visibility: 'public',
  enabled: true,
  sortOrder: 1,
  payload: { id: 'year_intro', type: 'intro', title: 'Your year' },
  when: { metric: 'orders_count', op: 'gte', value: 1 },
  createdAt: '2026-01-01T00:00:00.000Z',
  updatedAt: '2026-01-01T00:00:00.000Z',
};

describe('StoriesPage', () => {
  beforeEach(() => {
    navigate.mockReset();
    vi.mocked(useGetApiAdminStories).mockReturnValue({
      data: { items: [story] },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useGetApiAdminStories>);
  });

  it('renders stories and navigates on row click', async () => {
    const user = userEvent.setup();
    const { getByText, getByRole } = renderWithProviders(<StoriesPage />);

    expect(getByRole('link', { name: 'Create' })).toHaveAttribute('href', '/stories/new');
    expect(getByText('year_intro')).toBeInTheDocument();
    expect(getByRole('cell', { name: 'Intro' })).toBeInTheDocument();
    expect(getByRole('cell', { name: 'Public' })).toBeInTheDocument();
    expect(getByText('orders_count')).toBeInTheDocument();

    await user.click(getByText('year_intro'));

    expect(navigate).toHaveBeenCalledWith({ to: '/stories/$id', params: { id: 'year_intro' } });
  });

  it('reads filters from search params and writes them back', async () => {
    const user = userEvent.setup();
    const onUrlUpdate = vi.fn();
    const { getByLabelText } = renderWithProviders(<StoriesPage />, {
      searchParams: { search: 'year', enabled: 'true' },
      onUrlUpdate,
    });

    expect(getByLabelText('Search')).toHaveValue('year');

    await user.type(getByLabelText('Search'), 'x');

    const lastUpdate = onUrlUpdate.mock.calls.at(-1)?.[0] as
      | { searchParams: URLSearchParams }
      | undefined;
    expect(lastUpdate?.searchParams.get('search')).toBe('yearx');
    expect(lastUpdate?.searchParams.get('enabled')).toBe('true');
  });
});
