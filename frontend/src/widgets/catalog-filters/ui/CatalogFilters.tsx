import { NumberInput, Select, SimpleGrid, Stack, TextInput } from '@mantine/core';
import { parsePriceFilter, useProductFilters } from '@/features/product-filters';

const SORT_OPTIONS = [
  { value: 'title', label: 'По названию' },
  { value: 'price', label: 'По цене' },
  { value: 'rating', label: 'По рейтингу' },
];

const ORDER_OPTIONS = [
  { value: 'asc', label: 'По возрастанию' },
  { value: 'desc', label: 'По убыванию' },
];

type CatalogFiltersProps = {
  categories: string[];
};

export function CatalogFilters({ categories }: CatalogFiltersProps) {
  const [filters, setFilters] = useProductFilters();

  return (
    <Stack gap="sm">
      <SimpleGrid cols={{ base: 1, xs: 2, md: 4 }} spacing="sm">
        <TextInput
          label="Поиск"
          placeholder="Например, iPhone"
          value={filters.q}
          onChange={(event) => setFilters({ q: event.currentTarget.value, page: 1 })}
        />
        <Select
          label="Категория"
          placeholder="Все категории"
          clearable
          searchable
          data={categories.map((category) => ({ value: category, label: category }))}
          value={filters.category || null}
          onChange={(value) => setFilters({ category: value ?? '', page: 1 })}
        />
        <Select
          label="Сортировка"
          data={SORT_OPTIONS}
          value={filters.sort}
          onChange={(value) => setFilters({ sort: value ?? 'title', page: 1 })}
        />
        <Select
          label="Порядок"
          data={ORDER_OPTIONS}
          value={filters.order}
          onChange={(value) =>
            setFilters({ order: (value as 'asc' | 'desc' | null) ?? 'asc', page: 1 })
          }
        />
      </SimpleGrid>
      <SimpleGrid cols={{ base: 1, xs: 2 }} spacing="sm">
        <NumberInput
          label="Цена от"
          placeholder="0"
          min={0}
          value={parsePriceFilter(filters.priceMin) ?? undefined}
          onChange={(value) =>
            setFilters({
              priceMin: typeof value === 'number' ? String(value) : '',
              page: 1,
            })
          }
        />
        <NumberInput
          label="Цена до"
          placeholder="100000"
          min={0}
          value={parsePriceFilter(filters.priceMax) ?? undefined}
          onChange={(value) =>
            setFilters({
              priceMax: typeof value === 'number' ? String(value) : '',
              page: 1,
            })
          }
        />
      </SimpleGrid>
    </Stack>
  );
}
