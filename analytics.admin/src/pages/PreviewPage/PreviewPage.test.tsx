import { waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fetchReactPackument, getApiAdminPreview, loadRecapEngine } from '@/entities/preview';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import PreviewPage from './PreviewPage';

vi.mock('@/entities/preview', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/entities/preview')>();

  return {
    ...actual,
    fetchReactPackument: vi.fn(),
    loadRecapEngine: vi.fn(),
    getApiAdminPreview: vi.fn(),
  };
});

function FakeRecap({ autoplay }: { autoplay?: boolean }) {
  return <div>{autoplay ? 'Recap autoplay-on' : 'Recap autoplay-off'}</div>;
}

describe('PreviewPage', () => {
  beforeEach(() => {
    vi.mocked(fetchReactPackument).mockResolvedValue({
      'dist-tags': { latest: '2.0.1' },
      versions: {
        '2.0.1': { dependencies: { '@recap-engine/core': '^1.3.0' } },
        '1.0.0': { dependencies: { '@recap-engine/core': '1.0.0' } },
      },
    });
    vi.mocked(getApiAdminPreview).mockResolvedValue({ schemaVersion: 1 });
    vi.mocked(loadRecapEngine).mockResolvedValue({
      Recap: FakeRecap,
      prepareRecap: (payload) => ({ data: payload, scenes: [] }),
      createTheme: () => ({ cssVars: {} }),
    });
  });

  it('loads versions, runs preview, and toggles autoplay without a second fetch', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole, findByText } = renderWithProviders(<PreviewPage />);

    await waitFor(() => {
      expect(getByRole('textbox', { name: '@recap-engine/react' })).toHaveValue('2.0.1');
    });
    expect(getByLabelText('@recap-engine/core')).toHaveValue('^1.3.0');
    expect(getByLabelText('@recap-engine/core')).toBeDisabled();

    await user.click(getByRole('button', { name: 'Preview' }));

    expect(await findByText('Recap autoplay-off')).toBeInTheDocument();
    expect(getApiAdminPreview).toHaveBeenCalledOnce();

    await user.click(getByRole('switch', { name: 'Autoplay' }));

    expect(await findByText('Recap autoplay-on')).toBeInTheDocument();
    expect(getApiAdminPreview).toHaveBeenCalledOnce();
  });
});
