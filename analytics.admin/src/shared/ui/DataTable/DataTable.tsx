import { Table, Text } from '@mantine/core';
import type { ReactNode } from 'react';
import DataTableLoader from './DataTableLoader';

export type DataTableColumn<T> = {
  key: string;
  header: ReactNode;
  render: (row: T) => ReactNode;
};

type DataTableProps<T> = {
  columns: DataTableColumn<T>[];
  items: T[];
  onRowClick?: (row: T) => void;
  emptyMessage?: string;
  getRowId?: (row: T) => string;
  isLoading?: boolean;
  skeletonRowCount?: number;
};

export default function DataTable<T>({
  columns,
  items,
  onRowClick,
  emptyMessage = 'No items',
  getRowId,
  isLoading = false,
  skeletonRowCount = 8,
}: DataTableProps<T>) {
  if (isLoading) {
    return <DataTableLoader columns={columns} skeletonRowCount={skeletonRowCount} />;
  }

  if (items.length === 0) {
    return <Text c="dimmed">{emptyMessage}</Text>;
  }

  return (
    <Table highlightOnHover={Boolean(onRowClick)} striped>
      <Table.Thead>
        <Table.Tr>
          {columns.map((column) => (
            <Table.Th key={column.key}>{column.header}</Table.Th>
          ))}
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {items.map((item, index) => (
          <Table.Tr
            key={getRowId ? getRowId(item) : index}
            onClick={onRowClick ? () => onRowClick(item) : undefined}
            style={onRowClick ? { cursor: 'pointer' } : undefined}
          >
            {columns.map((column) => (
              <Table.Td key={column.key}>{column.render(item)}</Table.Td>
            ))}
          </Table.Tr>
        ))}
      </Table.Tbody>
    </Table>
  );
}
