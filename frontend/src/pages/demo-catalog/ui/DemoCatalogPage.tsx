import { Alert, Center, Loader, Pagination, Stack } from '@mantine/core';
import { getDemoUserById } from '@/entities/demo-user';
import { useOpenRecapModal } from '@/features/open-recap-modal';
import { useProductFilters } from '@/features/product-filters';
import { env } from '@/shared/config';
import { AvitoHeader } from '@/widgets/avito-header';
import { CatalogFilters } from '@/widgets/catalog-filters';
import { ProductGrid } from '@/widgets/product-grid';
import { RecapModal } from '@/widgets/recap-modal';
import { RecapYearBanner } from '@/widgets/recap-year-banner';
import { useProductsCatalog } from '../model/use-products-catalog';

type DemoCatalogPageProps = {
  userId: string;
};

export function DemoCatalogPage({ userId }: DemoCatalogPageProps) {
  const user = getDemoUserById(userId);
  const { isOpen, open, close } = useOpenRecapModal();
  const [filters, setFilters] = useProductFilters();
  const { products, categories, total, isLoading, isError, error } = useProductsCatalog();

  const totalPages = Math.max(1, Math.ceil(total / 12));

  return (
    <div className="page-shell">
      <AvitoHeader
        userId={userId}
        searchValue={filters.q}
        onSearchChange={(value) => setFilters({ q: value, page: 1 })}
        onSearchSubmit={() => setFilters({ page: 1 })}
      />

      <main className="page-content">
        <RecapYearBanner year={env.recapYear} onOpen={open} />
        <CatalogFilters categories={categories} />

        {isLoading ? (
          <Center py="xl">
            <Loader color="avito" />
          </Center>
        ) : null}

        {isError ? (
          <Alert color="red" title="Не удалось загрузить объявления">
            {error instanceof Error ? error.message : 'Unknown error'}
          </Alert>
        ) : null}

        {!isLoading && !isError ? (
          <Stack gap="md">
            <ProductGrid
              title={user ? `Рекомендации для ${user.name}` : 'Рекомендации для вас'}
              products={products}
            />
            <Pagination
              value={filters.page}
              onChange={(page) => setFilters({ page })}
              total={totalPages}
              color="avito"
            />
          </Stack>
        ) : null}
      </main>

      <RecapModal userId={userId} opened={isOpen} onClose={close} />
    </div>
  );
}
