import {
  closestCenter,
  DndContext,
  type DragEndEvent,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Table, Text, UnstyledButton } from '@mantine/core';
import type { CSSProperties, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import DataTableLoader from './DataTableLoader';

export type DataTableColumn<T> = {
  key: string;
  header: ReactNode;
  render: (row: T) => ReactNode;
};

type DataTableBaseProps<T> = {
  columns: DataTableColumn<T>[];
  items: T[];
  onRowClick?: (row: T) => void;
  emptyMessage?: string;
  isLoading?: boolean;
  skeletonRowCount?: number;
};

type DataTableProps<T> = DataTableBaseProps<T> &
  (
    | {
        onReorder: (items: T[], movedId: string) => void;
        getRowId: (row: T) => string;
      }
    | {
        onReorder?: undefined;
        getRowId?: (row: T) => string;
      }
  );

type SortableRowProps<T> = {
  item: T;
  id: string;
  columns: DataTableColumn<T>[];
  onRowClick?: (row: T) => void;
  dragHandleLabel: string;
};

function SortableRow<T>({ item, id, columns, onRowClick, dragHandleLabel }: SortableRowProps<T>) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id,
  });

  const style: CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    cursor: onRowClick ? 'pointer' : undefined,
    opacity: isDragging ? 0.8 : undefined,
  };

  return (
    <Table.Tr
      ref={setNodeRef}
      onClick={onRowClick ? () => onRowClick(item) : undefined}
      style={style}
    >
      <Table.Td>
        <UnstyledButton
          {...attributes}
          {...listeners}
          aria-label={dragHandleLabel}
          onClick={(event) => event.stopPropagation()}
          style={{ cursor: 'grab', touchAction: 'none' }}
        >
          ⠿
        </UnstyledButton>
      </Table.Td>
      {columns.map((column) => (
        <Table.Td key={column.key}>{column.render(item)}</Table.Td>
      ))}
    </Table.Tr>
  );
}

export default function DataTable<T>({
  columns,
  items,
  onRowClick,
  emptyMessage = 'No items',
  getRowId,
  isLoading = false,
  skeletonRowCount = 8,
  onReorder,
}: DataTableProps<T>) {
  const { t } = useTranslation();
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  const loaderColumns = onReorder ? [{ key: 'drag', header: '' }, ...columns] : columns;

  if (isLoading) {
    return <DataTableLoader columns={loaderColumns} skeletonRowCount={skeletonRowCount} />;
  }

  if (items.length === 0) {
    return <Text c="dimmed">{emptyMessage}</Text>;
  }

  const rowIds = getRowId ? items.map(getRowId) : [];

  function handleDragEnd(event: DragEndEvent) {
    if (!onReorder || !getRowId) {
      return;
    }

    const { active, over } = event;
    if (!over || active.id === over.id) {
      return;
    }

    const from = rowIds.indexOf(String(active.id));
    const to = rowIds.indexOf(String(over.id));
    if (from < 0 || to < 0) {
      return;
    }

    onReorder(arrayMove(items, from, to), String(active.id));
  }

  const headerRow = (
    <Table.Thead>
      <Table.Tr>
        {onReorder ? <Table.Th /> : null}
        {columns.map((column) => (
          <Table.Th key={column.key}>{column.header}</Table.Th>
        ))}
      </Table.Tr>
    </Table.Thead>
  );

  if (!onReorder || !getRowId) {
    return (
      <Table highlightOnHover={Boolean(onRowClick)} striped>
        {headerRow}
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

  const dragHandleLabel = t('table.dragHandle');

  return (
    <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
      <Table highlightOnHover={Boolean(onRowClick)} striped>
        {headerRow}
        <Table.Tbody>
          <SortableContext items={rowIds} strategy={verticalListSortingStrategy}>
            {items.map((item) => {
              const id = getRowId(item);
              return (
                <SortableRow
                  key={id}
                  id={id}
                  item={item}
                  columns={columns}
                  onRowClick={onRowClick}
                  dragHandleLabel={dragHandleLabel}
                />
              );
            })}
          </SortableContext>
        </Table.Tbody>
      </Table>
    </DndContext>
  );
}
