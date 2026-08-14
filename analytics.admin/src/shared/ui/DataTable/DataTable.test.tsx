import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import DataTable from './DataTable';

const columns = [{ key: 'title', header: 'Title', render: (row: { title: string }) => row.title }];

describe('DataTable', () => {
  it('renders rows', () => {
    const { getByText } = renderWithProviders(
      <DataTable columns={columns} items={[{ title: 'First' }, { title: 'Second' }]} />,
    );

    expect(getByText('First')).toBeInTheDocument();
    expect(getByText('Second')).toBeInTheDocument();
  });

  it('renders the empty state', () => {
    const { getByText } = renderWithProviders(
      <DataTable columns={columns} items={[]} emptyMessage="Nothing here" />,
    );

    expect(getByText('Nothing here')).toBeInTheDocument();
  });

  it('calls onRowClick with the row', async () => {
    const user = userEvent.setup();
    const onRowClick = vi.fn();
    const { getByText } = renderWithProviders(
      <DataTable columns={columns} items={[{ title: 'First' }]} onRowClick={onRowClick} />,
    );

    await user.click(getByText('First'));

    expect(onRowClick).toHaveBeenCalledWith({ title: 'First' });
  });

  it('renders a drag handle when onReorder is set', () => {
    const { getByLabelText } = renderWithProviders(
      <DataTable
        columns={columns}
        items={[{ title: 'First' }]}
        getRowId={(row) => row.title}
        onReorder={vi.fn()}
      />,
    );

    expect(getByLabelText('Reorder')).toBeInTheDocument();
  });

  it('still calls onRowClick when a reorderable row is clicked', async () => {
    const user = userEvent.setup();
    const onRowClick = vi.fn();
    const { getByText } = renderWithProviders(
      <DataTable
        columns={columns}
        items={[{ title: 'First' }]}
        getRowId={(row) => row.title}
        onReorder={vi.fn()}
        onRowClick={onRowClick}
      />,
    );

    await user.click(getByText('First'));

    expect(onRowClick).toHaveBeenCalledWith({ title: 'First' });
  });

  it('renders a table skeleton while loading', () => {
    const { getByLabelText, queryByText, getAllByRole } = renderWithProviders(
      <DataTable columns={columns} items={[]} isLoading skeletonRowCount={3} />,
    );

    expect(getByLabelText('Loading')).toBeInTheDocument();
    expect(queryByText('No items')).not.toBeInTheDocument();
    expect(getAllByRole('row')).toHaveLength(4);
  });
});
