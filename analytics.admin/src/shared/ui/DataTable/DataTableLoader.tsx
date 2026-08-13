import { Skeleton, Table } from '@mantine/core';
import type { ReactNode } from 'react';

type DataTableLoaderColumn = {
  key: string;
  header: ReactNode;
};

type DataTableLoaderProps = {
  columns: DataTableLoaderColumn[];
  skeletonRowCount?: number;
};

export default function DataTableLoader({ columns, skeletonRowCount = 8 }: DataTableLoaderProps) {
  const skeletonRows = Array.from(
    { length: skeletonRowCount },
    (_, rowIndex) => `skeleton-${rowIndex}`,
  );

  return (
    <Table aria-busy="true" aria-label="Loading">
      <Table.Thead>
        <Table.Tr>
          {columns.map((column) => (
            <Table.Th key={column.key}>{column.header}</Table.Th>
          ))}
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {skeletonRows.map((rowId) => (
          <Table.Tr key={rowId}>
            {columns.map((column) => (
              <Table.Td key={column.key}>
                <Skeleton height={22} radius="xl" />
              </Table.Td>
            ))}
          </Table.Tr>
        ))}
      </Table.Tbody>
    </Table>
  );
}
