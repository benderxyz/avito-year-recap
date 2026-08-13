import type { ReactNode } from 'react';
import { describe, expect, it, vi } from 'vitest';
import AppNav from '@/app/layout/AppNav';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';

vi.mock('@tanstack/react-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@tanstack/react-router')>();

  return {
    ...actual,
    useRouterState: ({
      select,
    }: {
      select: (state: { location: { pathname: string } }) => string;
    }) => select({ location: { pathname: '/metrics' } }),
    Link: ({ to, children, ...props }: { to: string; children?: ReactNode }) => (
      <a href={to} {...props}>
        {children}
      </a>
    ),
  };
});

describe('AppNav', () => {
  it('renders links to admin sections', () => {
    const { getByRole } = renderWithProviders(<AppNav />);

    expect(getByRole('link', { name: 'Metrics' })).toHaveAttribute('href', '/metrics');
    expect(getByRole('link', { name: 'Badges' })).toHaveAttribute('href', '/badges');
    expect(getByRole('link', { name: 'Stories' })).toHaveAttribute('href', '/stories');
    expect(getByRole('link', { name: 'Recommendations' })).toHaveAttribute(
      'href',
      '/recommendations',
    );
    expect(getByRole('link', { name: 'Preview' })).toHaveAttribute('href', '/preview');
  });
});
