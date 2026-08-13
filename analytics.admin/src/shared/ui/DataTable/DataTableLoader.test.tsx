import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import DataTableLoader from './DataTableLoader';

const columns = [{ key: 'title', header: 'Title' }];

describe('DataTableLoader', () => {
  it('renders header and skeleton rows', () => {
    const { getByLabelText, getByText, getAllByRole } = renderWithProviders(
      <DataTableLoader columns={columns} skeletonRowCount={3} />,
    );

    expect(getByLabelText('Loading')).toBeInTheDocument();
    expect(getByText('Title')).toBeInTheDocument();
    expect(getAllByRole('row')).toHaveLength(4);
  });
});
