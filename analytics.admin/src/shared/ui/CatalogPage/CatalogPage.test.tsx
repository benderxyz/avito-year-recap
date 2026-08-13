import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import CatalogPage from './CatalogPage';

describe('CatalogPage', () => {
  it('renders title, filters and children', () => {
    const { getByText } = renderWithProviders(
      <CatalogPage title="Badges" filters={<div>Filters</div>}>
        <div>Table</div>
      </CatalogPage>,
    );

    expect(getByText('Badges')).toBeInTheDocument();
    expect(getByText('Filters')).toBeInTheDocument();
    expect(getByText('Table')).toBeInTheDocument();
  });

  it('shows an error alert', () => {
    const { getByText } = renderWithProviders(
      <CatalogPage title="Badges" isError errorMessage="Nope">
        <div>Table</div>
      </CatalogPage>,
    );

    expect(getByText('Nope')).toBeInTheDocument();
  });

  it('renders actions next to the title', () => {
    const { getByRole } = renderWithProviders(
      <CatalogPage title="Metrics" actions={<button type="button">Create</button>}>
        <div>Table</div>
      </CatalogPage>,
    );

    expect(getByRole('button', { name: 'Create' })).toBeInTheDocument();
  });
});
